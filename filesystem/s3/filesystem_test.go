package s3

import (
	"github.com/goamz/goamz/aws"
	//"github.com/goamz/goamz/s3"
	"github.com/mefellows/mirror/filesystem"
	"os"
	"testing"
)

var oldAuthFunc = auth

func dummyAuth() {

	auth = func() (*aws.Auth, error) {
		return &aws.Auth{}, nil
	}
}

func restoreAuth() {
	auth = oldAuthFunc
}

func TestNew(t *testing.T) {
	dummyAuth()
	s3, err := New("s3://mybucket.s3.amazonaws.com")
	if err != nil {
		t.Fatalf("Got error %s", err.Error())
	}
	if s3.config.bucket != "mybucket" {
		t.Fatalf("Expected bucket to be 'mybucket', got %s", s3.config.bucket)
	}
	restoreAuth()
}

func TestNew_AuthWithEnvironment(t *testing.T) {
	_, err := New("s3://mybucket.s3.amazonaws.com")
	os.Setenv("AWS_ACCESS_KEY_ID", "bar")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "bat")
	if err != nil {
		t.Fatalf("Got error %s", err.Error())
	}
}

func TestNew_AuthWithFileError(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	os.Setenv("AWS_PROFILE", "mirrors3testprofile")
	_, err := New("s3://mybucket.s3.amazonaws.com")
	if err == nil {
		t.Fatalf("Expected error. Note that if this test fails, it's likely there exists a `$HOME/.aws/credentials` file that's available for the auth() method, with the profile `mirror3testprofile` - weird.")
	}
}

func TestConfig_VirtualhostStandard(t *testing.T) {

	// US Standard: us-east-1

	// Virtualhost bucket
	conf, err := config("s3://mybucket.s3.amazonaws.com/foo/bar.txt")
	if err != nil {
		t.Fatalf("Got unexpected error %s", err.Error())
	}
	if conf.baseURL != "s3://mybucket.s3.amazonaws.com" {
		t.Fatalf("Base URL name should be 's3://mybucket.s3.amazonaws.com' but was %s", conf.baseURL)
	}
	if conf.bucket != "mybucket" {
		t.Fatalf("Bucket name should be 'mybucket' but was %s", conf.bucket)
	}
	if conf.region != "us-east-1" {
		t.Fatalf("Region should be 'us-east-1' but was %s", conf.region)
	}

}
func TestConfig_VirtualhostOtherRegion(t *testing.T) {

	conf, err := config("s3://mybucket.s3-us-west-1.amazonaws.com/foo/bar.txt")

	if err != nil {
		t.Fatalf("Got unexpected error %s", err.Error())
	}
	if conf.baseURL != "s3://mybucket.s3-us-west-1.amazonaws.com" {
		t.Fatalf("Base URL name should be 's3://mybucket.s3-us-west-1.amazonaws.com' but was %s", conf.baseURL)
	}
	if conf.bucket != "mybucket" {
		t.Fatalf("Bucket name should be 'mybucket' but was %s", conf.bucket)
	}
	if conf.region != "us-west-1" {
		t.Fatalf("Region should be 'us-west-1' but was %s", conf.region)
	}
}
func TestConfig_PathStyleStandard(t *testing.T) {

	conf, err := config("s3://s3.amazonaws.com/mybucket/foo/bar.txt")
	if err != nil {
		t.Fatalf("Got unexpected error %s", err.Error())
	}
	if conf.baseURL != "s3://s3.amazonaws.com/mybucket" {
		t.Fatalf("Base URL name should be 's3://s3.amazonaws.com/mybucket' but was %s", conf.baseURL)
	}
	if conf.bucket != "mybucket" {
		t.Fatalf("Bucket name should be 'mybucket' but was %s", conf.bucket)
	}
	if conf.region != "us-east-1" {
		t.Fatalf("Region should be 'us-east-1' but was %s", conf.region)
	}
}
func TestConfig_PathStyleOtherRegion(t *testing.T) {

	conf, err := config("s3://s3-us-west-1.amazonaws.com/mybucket/foo/bar.txt")
	if err != nil {
		t.Fatalf("Got unexpected error %s", err.Error())
	}
	if conf.baseURL != "s3://s3-us-west-1.amazonaws.com/mybucket" {
		t.Fatalf("Base URL name should be 's3://s3-us-west-1.amazonaws.com/mybucket' but was %s", conf.baseURL)
	}
	if conf.bucket != "mybucket" {
		t.Fatalf("Bucket name should be 'mybucket' but was %s", conf.bucket)
	}
	if conf.region != "us-west-1" {
		t.Fatalf("Region should be 'us-west-1' but was %s", conf.region)
	}
}
func TestConfig_InvalidURL(t *testing.T) {

	_, err := config("s3://s3-us-west-1.amazonaws.com/mybucket")
	if err == nil {
		t.Fatalf("Expected error")
	}
	_, err = config("s3://s3-us-west-1.amazonaws.com/")
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestExt(t *testing.T) {
	file := filesystem.File{FileName: "/foo/bar/baz.txt"}
	if ext(file) != ".txt" {
		t.Fatalf("Expected .txt extension, got %s", ext(file))
	}
}

func TestMime(t *testing.T) {
	file := filesystem.File{FileName: "/foo/bar/baz.txt"}
	if mimeType(file) != "text/plain; charset=utf-8" {
		t.Fatalf("Expected text/plain mime, got %s", mimeType(file))
	}

	file = filesystem.File{FileName: "baz.json"}
	if mimeType(file) != "application/json" {
		t.Fatalf("Expected application/json mime, got %s", mimeType(file))
	}
}
