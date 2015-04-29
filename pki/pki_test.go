package pki

import (
	"github.com/mefellows/mirror/mirror"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var (
	tmpDir, _ = ioutil.TempDir("", "machine-test-")
)

func defaultPki() *PKI {
	os.Setenv("MIRROR_HOME", tmpDir)
	pki, _ := New()
	return pki
}

// Sets up a fake CA cert in our temp location. NOTE: It is the callers'
// responsibilty to issue a `defer os.RemoveAll(tmpDir)` once done
func generateCaCert() error {
	// Setup fake cert
	os.Mkdir(mirror.GetCADir(), 755)
	os.Mkdir(mirror.GetCertDir(), 755)
	caCertPath := filepath.Join(mirror.GetCADir(), "ca.pem")
	caKeyPath := filepath.Join(mirror.GetCertDir(), "key.pem")
	testOrg := "test-org"
	bits := 2048
	if err := GenerateCACertificate(caCertPath, caKeyPath, testOrg, bits); err != nil {
		return err
	}

	if _, err := os.Stat(caCertPath); err != nil {
		return err
	}
	if _, err := os.Stat(caKeyPath); err != nil {
		return err
	}
	return nil
}

func TestDefaultConfig(t *testing.T) {
	pki := defaultPki()

	expectedCaCertPath := path.Join(tmpDir, "/ca/ca.pem")
	if pki.Config.caCertPath != expectedCaCertPath {
		t.Fatalf("Expected CA Cert path to be %s, but got %s", expectedCaCertPath, pki.Config.caCertPath)
	}

	expectedCaKeyPath := path.Join(tmpDir, "/ca/key.pem")
	if pki.Config.caKeyPath != expectedCaKeyPath {
		t.Fatalf("Expected CA Key path to be %s, but got %s", expectedCaKeyPath, pki.Config.caKeyPath)
	}
}

func TestDiscoverCAs(t *testing.T) {
	// cleanup
	generateCaCert()
	defer os.RemoveAll(tmpDir)

	pki := defaultPki()
	pool, err := pki.discoverCAs()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(pool.Subjects()) == 0 {
		t.Fatalf("Empty cert pool!")
	}

}

func TestSetupPKI(t *testing.T) {

}

func TestRemovePKI(t *testing.T) {

}
