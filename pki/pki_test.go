// General notes on this package:
//   1. Needs all non-PKI specific stuff removed
//   2. Needs command related logging removed
//   3. Ideally this package can be extracted from this project and made more generic
package pki

import (
	"crypto/tls"
	"github.com/mefellows/mirror/mirror"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	tmpDir, _ = ioutil.TempDir("", "machine-test-")
)

func defaultPki() *PKI {
	os.Setenv("MIRROR_HOME", tmpDir)
	pki, _ := New()
	return pki
}

func TestNew(t *testing.T) {
	os.Setenv("MIRROR_HOME", tmpDir)
	pki, err := New()

	if pki.Config.Insecure == true {
		t.Fatalf("PKI default should be secure")
	}

	if err != nil {
		t.Fatalf("Did not expect error:  %s", err)
	}
	if pki == nil {
		t.Fatalf("Expected PKI to return non-nil object")
	}

	paths := []string{pki.Config.CaCertPath, pki.Config.CaKeyPath, pki.Config.ClientCertPath, pki.Config.ClientKeyPath, pki.Config.ServerKeyPath, pki.Config.ServerCertPath}
	for _, file := range paths {
		if _, err := os.Stat(file); err != nil {
			time.Sleep(time.Second * 30)
		}
	}

	os.RemoveAll(tmpDir)
}

func TestNewWithConfig(t *testing.T) {
	os.Setenv("MIRROR_HOME", tmpDir)
	config := &Config{
		Insecure:       true,
		CaCertPath:     path.Join(tmpDir, "ca", "ca.pem"),
		CaKeyPath:      path.Join(tmpDir, "ca", "key.pem"),
		ClientCertPath: path.Join(tmpDir, "certs", "cert.pem"),
		ClientKeyPath:  path.Join(tmpDir, "certs", "cert-key.pem"),
		ServerCertPath: path.Join(tmpDir, "certs", "server-cert.pem"),
		ServerKeyPath:  path.Join(tmpDir, "certs", "server-key.pem"),
	}
	pki, err := NewWithConfig(config)

	if err != nil {
		t.Fatalf("Did not expect error:  %s", err)
	}
	if pki == nil {
		t.Fatalf("Expected PKI to return non-nil object")
	}

	paths := []string{config.CaCertPath, config.CaKeyPath, config.ClientCertPath, config.ClientKeyPath, config.ServerKeyPath, config.ServerCertPath}
	for _, file := range paths {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("File '%s' did not exist", file)
		}
	}

	os.RemoveAll(tmpDir)
}

func TestRemoveAll(t *testing.T) {
	os.Setenv("MIRROR_HOME", tmpDir)
	pki, err := New()

	if pki.Config.Insecure == true {
		t.Fatalf("PKI default should be secure")
	}

	if err != nil {
		t.Fatalf("Did not expect error:  %s", err)
	}
	if pki == nil {
		t.Fatalf("Expected PKI to return non-nil object")
	}

	err = pki.RemovePKI()

	paths := []string{pki.Config.CaCertPath, pki.Config.CaKeyPath, pki.Config.ClientCertPath, pki.Config.ClientKeyPath, pki.Config.ServerKeyPath, pki.Config.ServerCertPath}
	for _, file := range paths {
		if _, err := os.Stat(file); err == nil {
			t.Fatalf("File '%s' still exists, but should have been removed", file)
		}
	}
}

// Sets up a fake CA cert in our temp location. NOTE: It is the callers'
// responsibilty to issue a `defer os.RemoveAll(tmpDir)` once done
func generateCaCert() error {
	// Setup fake cert
	os.Mkdir(mirror.GetCADir(), 755)
	os.Mkdir(mirror.GetCertDir(), 755)
	CaCertPath := filepath.Join(mirror.GetCADir(), "ca.pem")
	CaKeyPath := filepath.Join(mirror.GetCertDir(), "key.pem")
	testOrg := "test-org"
	bits := 2048
	if err := GenerateCACertificate(CaCertPath, CaKeyPath, testOrg, bits); err != nil {
		return err
	}

	if _, err := os.Stat(CaCertPath); err != nil {
		return err
	}
	if _, err := os.Stat(CaKeyPath); err != nil {
		return err
	}
	return nil
}

func TestDefaultConfig(t *testing.T) {
	pki := defaultPki()

	expectedCaCertPath := path.Join(tmpDir, "/ca/ca.pem")
	if pki.Config.CaCertPath != expectedCaCertPath {
		t.Fatalf("Expected CA Cert path to be %s, but got %s", expectedCaCertPath, pki.Config.CaCertPath)
	}

	expectedCaKeyPath := path.Join(tmpDir, "/ca/key.pem")
	if pki.Config.CaKeyPath != expectedCaKeyPath {
		t.Fatalf("Expected CA Key path to be %s, but got %s", expectedCaKeyPath, pki.Config.CaKeyPath)
	}
	os.RemoveAll(tmpDir)
}

func TestDiscoverCAs(t *testing.T) {
	generateCaCert()

	pki := defaultPki()
	pool, err := pki.discoverCAs()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(pool.Subjects()) == 0 {
		t.Fatalf("Empty cert pool!")
	}
	if len(pool.Subjects()) != 1 {
		t.Fatalf("More subjects than the (1) expected, got %d", len(pool.Subjects()))
	}

	// Manually add extra CAs and check they are imported
	cert, _ := ioutil.ReadFile(pki.Config.CaCertPath)
	key, _ := ioutil.ReadFile(pki.Config.CaKeyPath)
	ioutil.WriteFile(filepath.Join(filepath.Dir(pki.Config.CaCertPath), "ca-test.pem"), cert, 0600)
	ioutil.WriteFile(filepath.Join(filepath.Dir(pki.Config.CaCertPath), "key-test.pem"), key, 0600)
	generateCaCert()

	pool, err = pki.discoverCAs()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(pool.Subjects()) == 0 {
		t.Fatalf("Empty cert pool!")
	}
	if len(pool.Subjects()) != 2 {
		t.Fatalf("More subjects than the (2) expected, got %d", len(pool.Subjects()))
	}

	// TODO: Check that certificates created against them are valid?
	os.RemoveAll(tmpDir)

}

func TestOutputKeysAndThings(t *testing.T) {
	pki := defaultPki()
	output, _ := pki.OutputCAKey()
	if !strings.Contains(output, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Fatalf("Expected to see \"-----BEGIN RSA PRIVATE KEY-----\", but got \"%s\"", output)
	}

	output, _ = pki.OutputCACert()
	if !strings.Contains(output, "-----BEGIN CERTIFICATE-----") {
		t.Fatalf("Expected to see \"-----BEGIN CERTIFICATE-----\", but got \"%s\"", output)

	}

	output, _ = pki.OutputClientKey()
	if !strings.Contains(output, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Fatalf("Expected to see \"-----BEGIN RSA PRIVATE KEY-----\", but got \"%s\"", output)
	}

	output, _ = pki.OutputClientCert()
	if !strings.Contains(output, "-----BEGIN CERTIFICATE-----") {
		t.Fatalf("Expected to see \"-----BEGIN CERTIFICATE-----\", but got \"%s\"", output)

	}
	os.RemoveAll(tmpDir)
}

func TestGetServerTLSConfig(t *testing.T) {
	pki := defaultPki()
	config, _ := pki.GetServerTLSConfig()
	if config.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Fatalf("Communications should be secure by default, got: %s", config.ClientAuth)
	}

	os.Setenv("MIRROR_HOME", tmpDir)
	pkiConfig := &Config{
		Insecure:       true,
		CaCertPath:     path.Join(tmpDir, "ca", "ca.pem"),
		CaKeyPath:      path.Join(tmpDir, "ca", "key.pem"),
		ClientCertPath: path.Join(tmpDir, "certs", "cert.pem"),
		ClientKeyPath:  path.Join(tmpDir, "certs", "cert-key.pem"),
		ServerCertPath: path.Join(tmpDir, "certs", "server-cert.pem"),
		ServerKeyPath:  path.Join(tmpDir, "certs", "server-key.pem"),
	}
	pki, _ = NewWithConfig(pkiConfig)
	config, _ = pki.GetServerTLSConfig()
	if config.ClientAuth != tls.NoClientCert {
		t.Fatalf("Secure communications disabled, should not check client cert (tls.NoClientCert) but instead got: %s", config.ClientAuth)
	}

	// Delete the CA - This should actually be an error as we need a non-nil certPool
	os.RemoveAll(tmpDir)
	config, err := pki.GetServerTLSConfig()
	if err == nil {
		t.Fatalf("no CA/Server Certs, even in --insecure mode this should cause an issue due to TLS library requirements for a CertPool. Happy days if not.")
	}

	// Delete the CA - we should get an error
	pkiConfig.Insecure = false
	config, err = pki.GetServerTLSConfig()
	if err == nil {
		t.Fatalf("No CA present, should be an error")
	}
}

func TestGetClientTLSConfig(t *testing.T) {
	pki := defaultPki()
	config, _ := pki.GetClientTLSConfig()
	if config.InsecureSkipVerify != false {
		t.Fatalf("Communications should be secure by default, got: %s", config.ClientAuth)
	}

	os.Setenv("MIRROR_HOME", tmpDir)
	pkiConfig := &Config{
		Insecure:       true,
		CaCertPath:     path.Join(tmpDir, "ca", "ca.pem"),
		CaKeyPath:      path.Join(tmpDir, "ca", "key.pem"),
		ClientCertPath: path.Join(tmpDir, "certs", "cert.pem"),
		ClientKeyPath:  path.Join(tmpDir, "certs", "cert-key.pem"),
		ServerCertPath: path.Join(tmpDir, "certs", "server-cert.pem"),
		ServerKeyPath:  path.Join(tmpDir, "certs", "server-key.pem"),
	}

	// Insecure
	pki, _ = NewWithConfig(pkiConfig)
	config, _ = pki.GetClientTLSConfig()
	if config.InsecureSkipVerify == false {
		t.Fatalf("Secure communications disabled, got: %s", config.InsecureSkipVerify)
	}

	// Insecure - no client certs also. Should not error
	os.Remove(path.Join(tmpDir, "certs", "cert.pem"))
	os.Remove(path.Join(tmpDir, "certs", "cert-key.pem"))
	pki, _ = NewWithConfig(pkiConfig)
	config, err := pki.GetClientTLSConfig()
	if err != nil {
		t.Fatalf("Did not expect err: %s", err)
	}
	if config.InsecureSkipVerify == false {
		t.Fatalf("Secure communications disabled, got: %s", config.InsecureSkipVerify)
	}

	pkiConfig.Insecure = false
	pki, _ = NewWithConfig(pkiConfig)
	config, err = pki.GetClientTLSConfig()
	if err == nil {
		t.Fatalf("Expected error but did not get one")
	}
}

func TestGenerateClientCertificate(t *testing.T) {

}

func TestDiscoverClientCertificates(t *testing.T) {

}

func TestImportCA(t *testing.T) {

	// Happy scenario

	pki := defaultPki()
	pool, err := pki.discoverCAs()
	crt, _ := ioutil.ReadFile(pki.Config.CaCertPath)
	newPem := path.Join(tmpDir, "ca", "ca-import.pem")
	ioutil.WriteFile(newPem, crt, 0600)

	err = pki.ImportCA("mynewca", newPem)
	if err != nil {
		t.Fatalf("Did not expect error")
	}
	if len(pool.Subjects()) == 0 {
		t.Fatal("CA import failed")
	}

	os.RemoveAll(tmpDir)

	// Validate file doesn't exist
	pki = defaultPki()
	os.RemoveAll(tmpDir)
	pool, err = pki.discoverCAs()
	newPem = path.Join(tmpDir, "ca", "ca-import.pem")
	ioutil.WriteFile(newPem, crt, 0600)

	err = pki.ImportCA("mynewca", newPem)
	if err == nil {
		t.Fatalf("Expected error")
	}
	if len(pool.Subjects()) != 0 {
		t.Fatal("CA import should have faild, removed all certs")
	}

	// Validate name regex
	pki = defaultPki()
	newPem = path.Join(tmpDir, "ca", "ca-import.pem")
	ioutil.WriteFile(newPem, crt, 0600)
	names := []string{"&invalid", "in valid", "In^alid"}
	for _, name := range names {
		err = pki.ImportCA(name, newPem)
		if err == nil {
			t.Fatalf("Expected error")
		}
	}
	names = []string{"valid", "val1d", "val-1d-0123456789", "VAL1d", "val_id", "valid.fr33"}
	for _, name := range names {
		err = pki.ImportCA(name, newPem)
		if err != nil {
			t.Fatalf("Did not expect error: %s", err.Error())
		}
	}

	// Validate invalid CA
	ioutil.WriteFile(newPem, []byte{}, 0600)
	err = pki.ImportCA("mynewca", newPem)
	if err == nil && !strings.HasPrefix(err.Error(), "Certificate provided is not valid") {
		t.Fatalf("Expected error to start with 'Certificate provided is not valid'")
	}
}
func TestImportClient(t *testing.T) {

	// Happy scenario

	pki := defaultPki()
	crt, _ := ioutil.ReadFile(pki.Config.ClientCertPath)
	key, _ := ioutil.ReadFile(pki.Config.ClientKeyPath)
	os.Remove(path.Join(tmpDir, "certs", "cert.pem"))
	os.Remove(path.Join(tmpDir, "certs", "cert-key.pem"))

	newCrt := path.Join(tmpDir, "certs", "client-cert.pem")
	newKey := path.Join(tmpDir, "certs", "client-key.pem")
	ioutil.WriteFile(newCrt, crt, 0600)
	ioutil.WriteFile(newKey, key, 0600)

	err := pki.ImportClientCertAndKey(newCrt, newKey)
	if err != nil {
		t.Fatalf("Did not expect error: %s", err)
	}

	// Check - do files exist?
	if _, err := os.Stat(path.Join(tmpDir, "certs", "client-cert.pem")); err != nil {
		t.Fatalf("Did not expect error: %s", err)
	}
	if _, err := os.Stat(path.Join(tmpDir, "certs", "client-key.pem")); err != nil {
		t.Fatalf("Did not expect error: %s", err)
	}

	// Validate invalid CA
	ioutil.WriteFile(newCrt, []byte{}, 0600)
	err = pki.ImportClientCertAndKey(newCrt, newKey)
	if err == nil && !strings.HasPrefix(err.Error(), "Certificate provided is not valid") {
		t.Fatalf("Expected error to start with 'Certificate provided is not valid'")
	}
}
