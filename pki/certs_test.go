// Adapted from the fine folks over at Docker: https://github.com/docker/machine/blob/382f71fdda6f689d4de465588e4a76e9b8ae2836/utils/certs.go
package pki

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCACertificate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "machine-test-")
	if err != nil {
		t.Fatal(err)
	}
	// cleanup
	defer os.RemoveAll(tmpDir)

	os.Setenv("MACHINE_DIR", tmpDir)
	caCertPath := filepath.Join(tmpDir, "ca.pem")
	caKeyPath := filepath.Join(tmpDir, "key.pem")
	testOrg := "test-org"
	bits := 2048
	if err := GenerateCACertificate(caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(caCertPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(caKeyPath); err != nil {
		t.Fatal(err)
	}
	os.Setenv("MACHINE_DIR", "")
}

func TestGenerateCertificate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "machine-test-")
	if err != nil {
		t.Fatal(err)
	}
	// cleanup
	defer os.RemoveAll(tmpDir)

	os.Setenv("MACHINE_DIR", tmpDir)
	caCertPath := filepath.Join(tmpDir, "ca.pem")
	caKeyPath := filepath.Join(tmpDir, "key.pem")
	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "cert-key.pem")
	testOrg := "test-org"
	bits := 2048
	if err := GenerateCACertificate(caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(caCertPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(caKeyPath); err != nil {
		t.Fatal(err)
	}
	os.Setenv("MACHINE_DIR", "")

	// Client Cert
	if err := GenerateCertificate([]string{""}, certPath, keyPath, caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("certificate not created at %s", certPath)
	}

	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("key not created at %s", keyPath)
	}

	// Hostname Cert
	if err := GenerateCertificate([]string{"foo.com"}, certPath, keyPath, caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("certificate not created at %s", certPath)
	}

	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("key not created at %s", keyPath)
	}

	// IP based cert
	if err := GenerateCertificate([]string{"127.0.0.1"}, certPath, keyPath, caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("certificate not created at %s", certPath)
	}

	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("key not created at %s", keyPath)
	}
}
