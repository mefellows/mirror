// General notes on this package:
//   1. Needs all non-PKI specific stuff removed
//   2. Needs command related logging removed
//   3. Ideally this package can be extracted from this project and made more generic
package pki

import (
	"crypto/tls"
	"fmt"
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

	paths := []string{pki.Config.caCertPath, pki.Config.caKeyPath, pki.Config.clientCertPath, pki.Config.clientKeyPath, pki.Config.serverKeyPath, pki.Config.serverCertPath}
	for _, file := range paths {
		if _, err := os.Stat(file); err != nil {
			fmt.Printf("file %s not found, take a look...")
			time.Sleep(time.Second * 30)
			t.Fatalf("File '%s' did not exist. Error: %s", file, err)
		}
	}

	os.RemoveAll(tmpDir)
}

func TestNewWithConfig(t *testing.T) {
	os.Setenv("MIRROR_HOME", tmpDir)
	config := &Config{
		Insecure:       true,
		caCertPath:     path.Join(tmpDir, "ca", "ca.pem"),
		caKeyPath:      path.Join(tmpDir, "ca", "key.pem"),
		clientCertPath: path.Join(tmpDir, "certs", "cert.pem"),
		clientKeyPath:  path.Join(tmpDir, "certs", "cert-key.pem"),
		serverCertPath: path.Join(tmpDir, "certs", "server-cert.pem"),
		serverKeyPath:  path.Join(tmpDir, "certs", "server-key.pem"),
	}
	pki, err := NewWithConfig(config)

	if err != nil {
		t.Fatalf("Did not expect error:  %s", err)
	}
	if pki == nil {
		t.Fatalf("Expected PKI to return non-nil object")
	}

	paths := []string{config.caCertPath, config.caKeyPath, config.clientCertPath, config.clientKeyPath, config.serverKeyPath, config.serverCertPath}
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

	paths := []string{pki.Config.caCertPath, pki.Config.caKeyPath, pki.Config.clientCertPath, pki.Config.clientKeyPath, pki.Config.serverKeyPath, pki.Config.serverCertPath}
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
	os.RemoveAll(tmpDir)
}

func TestDiscoverCAs(t *testing.T) {
	// cleanup
	generateCaCert()

	pki := defaultPki()
	pool, err := pki.discoverCAs()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(pool.Subjects()) == 0 {
		t.Fatalf("Empty cert pool!")
	}
	os.RemoveAll(tmpDir)

	// TODO: Manually add extra CAs and check they are imported

	// TODO: Check that certificates created against them are valid?

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
		caCertPath:     path.Join(tmpDir, "ca", "ca.pem"),
		caKeyPath:      path.Join(tmpDir, "ca", "key.pem"),
		clientCertPath: path.Join(tmpDir, "certs", "cert.pem"),
		clientKeyPath:  path.Join(tmpDir, "certs", "cert-key.pem"),
		serverCertPath: path.Join(tmpDir, "certs", "server-cert.pem"),
		serverKeyPath:  path.Join(tmpDir, "certs", "server-key.pem"),
	}
	pki, _ = NewWithConfig(pkiConfig)
	config, _ = pki.GetServerTLSConfig()
	if config.ClientAuth != tls.NoClientCert {
		t.Fatalf("Secure communications disabled, got: %s", config.ClientAuth)
	}
	// SSL
	// Insecure
	// Invalid SSL
	// Invalid Client Cert
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
		caCertPath:     path.Join(tmpDir, "ca", "ca.pem"),
		caKeyPath:      path.Join(tmpDir, "ca", "key.pem"),
		clientCertPath: path.Join(tmpDir, "certs", "cert.pem"),
		clientKeyPath:  path.Join(tmpDir, "certs", "cert-key.pem"),
		serverCertPath: path.Join(tmpDir, "certs", "server-cert.pem"),
		serverKeyPath:  path.Join(tmpDir, "certs", "server-key.pem"),
	}
	pki, _ = NewWithConfig(pkiConfig)
	config, _ = pki.GetClientTLSConfig()
	if config.InsecureSkipVerify == false {
		t.Fatalf("Secure communications disabled, got: %s", config.InsecureSkipVerify)
	}

	// SSL
	// Insecure
	// Invalid SSL
	// Invalid Client Cert
}
