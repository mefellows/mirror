package s3

import (
	"errors"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/mefellows/mirror/filesystem"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type WriteRequest struct {
	// TODO: How do we make this protocol agnostic??
	// "type not registered for interface: fs.StdFile"
	File filesystem.File
	//File fs.StdFile
	Data []byte
	Perm os.FileMode
}
type WriteResponse struct {
	Success bool
}

func (f *S3FileSystem) WriteRemote(req *WriteRequest, res *WriteResponse) error {
	return nil
}

// S3 File System implementation
type S3FileSystem struct {
	tree      filesystem.FileTree // Returns a FileTree structure of Files representing the FileSystem hierarchy
	auth      *aws.Auth
	s3service *s3.S3
	config    *S3Config
	url       string
	bucket    *s3.Bucket
}

type S3Config struct {
	bucket  string
	region  string
	baseURL string // The base component of the S3 URL e.g. s3://s3.amazonaws.com/mybucket/. This component should be removed in any PUTs
}

// Create a new S3FileSystem object. Requires an S3 URL to configure
func New(url string) (*S3FileSystem, error) {
	auth, err := auth()
	if err != nil {
		return nil, err
	}
	config, err := config(url)
	region := aws.Regions[config.region]
	service := s3.New(*auth, region)

	s3fs := &S3FileSystem{
		auth:      auth,
		s3service: service,
		url:       url,
		config:    config,
		bucket:    service.Bucket(config.bucket),
	}
	return s3fs, nil
}

// Get Authentication details from environment
var auth = func() (*aws.Auth, error) {
	// Check $HOME/.aws/credentials first
	auth, err := aws.SharedAuth()
	if err == nil {
		return &auth, err
	}

	// Check environment variables
	auth, err = aws.EnvAuth()
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// Extract the bucket name and region from an s3:// URL
func config(url string) (*S3Config, error) {
	var virtualhostMatch = regexp.MustCompile(`^(s3:\/\/([a-zA-Z-_\.0-9]+)\.s3\.amazonaws\.com)`)
	var virtualhostWithRegionMatch = regexp.MustCompile(`^(s3:\/\/([a-zA-Z-_\.0-9]+)\.s3-([a-zA-Z-_\.0-9]+)\.amazonaws\.com)`)
	var pathMatch = regexp.MustCompile(`^(s3:\/\/s3\.amazonaws\.com\/([a-zA-Z-_\.0-9]+))\/`)
	var pathWithRegionMatch = regexp.MustCompile(`^(s3:\/\/s3-([a-zA-Z-_\.0-9]+)\.amazonaws\.com\/([a-zA-Z-_\.0-9]+))\/`)
	var bucket string
	var baseURL string
	region := "us-east-1" // Default

	switch {
	case virtualhostMatch.MatchString(url):
		matches := virtualhostMatch.FindStringSubmatch(url)
		baseURL = matches[1]
		bucket = matches[2]
	case virtualhostWithRegionMatch.MatchString(url):
		matches := virtualhostWithRegionMatch.FindStringSubmatch(url)
		baseURL = matches[1]
		bucket = matches[2]
		region = matches[3]
	case pathMatch.MatchString(url):
		matches := pathMatch.FindStringSubmatch(url)
		baseURL = matches[1]
		bucket = matches[2]
	case pathWithRegionMatch.MatchString(url):
		matches := pathWithRegionMatch.FindStringSubmatch(url)
		baseURL = matches[1]
		region = matches[2]
		bucket = matches[3]
	default:
		return nil, errors.New("Invalid S3 URL provided")
	}

	return &S3Config{
		bucket:  bucket,
		region:  region,
		baseURL: baseURL,
	}, nil
}

func (fs S3FileSystem) Dir(dir string) ([]filesystem.File, error) {
	return nil, errors.New("Function not yet implemented")
}

func (fs S3FileSystem) Read(f filesystem.File) ([]byte, error) {
	return nil, errors.New("Function not yet implemented")
}

func (fs S3FileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) error {
	fileName := strings.TrimPrefix(file.Name(), fs.config.baseURL)
	return fs.bucket.Put(fileName, data, mimeType(file), s3.BucketOwnerFull, s3.Options{})
}

func (fs S3FileSystem) Delete(file filesystem.File) error {
	return errors.New("Function not yet implemented")
}

func (fs S3FileSystem) FileTree() filesystem.FileTree {
	return nil
}

type S3File struct {
	S3Name    string      // base name of the file
	S3Size    int64       // length in bytes for regular files; system-dependent for others
	S3ModTime time.Time   // modification time
	S3IsDir   bool        // abbreviation for Mode().IsDir()
	S3Mode    os.FileMode // abbreviation for Mode().IsDir()
}

func (f S3File) Name() string {
	return f.S3Name
}

func (f S3File) Size() int64 {
	return f.S3Size
}
func (f S3File) ModTime() time.Time {
	return f.S3ModTime
}
func (f S3File) IsDir() bool {
	return f.S3IsDir
}
func (f S3File) Mode() os.FileMode {
	return f.S3Mode
}
func (f S3File) Sys() interface{} {
	return nil
}
func ext(file filesystem.File) string {
	return filepath.Ext(file.Name())
}
func mimeType(file filesystem.File) string {
	return mime.TypeByExtension(ext(file))
}
