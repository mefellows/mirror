package pki

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
)

// Check for custom Certificate Authorities

// Check for server certificates (must be present for TLS to occur safely)

// Generate a client cert?

// General Mirror Pubic Key Infrastructure functions
type PKI struct {
}

type tlsConfig struct {
	clientKeyPath  string
	clientCertPath string
	serverKeyPath  string
	serverCertPath string
	caKeyPath      string
	caCertPath     string
}

func New() *PKI {
	return nil
}

func (p *PKI) RemovePKI() error {
	// TODO: Get Base Config containing Home Dir
	pkiHomeDir := "/Users/mfellows/.mirror.d/pki"
	err := os.RemoveAll(pkiHomeDir)

	return err
}

func (p *PKI) GenerateCert(hosts []string) (err error) {
	outputDir := filepath.Dir(".")
	pkiHomeDir := "/Users/mfellows/.mirror.d/pki"
	caCertPath := filepath.Join(pkiHomeDir, "ca.pem")
	caKeyPath := filepath.Join(pkiHomeDir, "key.pem")
	certPath := filepath.Join(outputDir, "cert.pem")
	keyPath := filepath.Join(outputDir, "cert-key.pem")
	testOrg := "client"
	bits := 2048

	if len(hosts) == 0 {
		hosts = []string{}
	}

	err = GenerateCert(hosts, certPath, keyPath, caCertPath, caKeyPath, testOrg, bits)
	if err == nil {
		_, err = os.Stat(certPath)
		_, err = os.Stat(keyPath)
	}
	return
}

// Validate all components of the PKI infrastructure are properly configured
func (p *PKI) CheckSetup() error {
	// Check directories

	// Check CA

	// Check server cert

	// Check client certs against CA (from conf + user?)

	// Check permissions?

	return nil
}

// Sets up the PKI infrastructure for client / server communications
// This involves creating directories, CAs, and client/server certs
func (p *PKI) SetupPKI(caHost string) error {

	// Ensure all paths are

	// TODO: Get Base Config containing Home Dir
	pkiHomeDir := "/Users/mfellows/.mirror.d/pki"
	os.Mkdir(pkiHomeDir, 0700)

	caCertPath := filepath.Join(pkiHomeDir, "ca.pem")
	caKeyPath := filepath.Join(pkiHomeDir, "key.pem")
	bits := 2048
	if _, err := os.Stat(caCertPath); err == nil {
		return fmt.Errorf("CA already exists. Run --delete to remove the old CA.")
	}

	if err := GenerateCACertificate(caCertPath, caKeyPath, caHost, bits); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(caCertPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(caKeyPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	serverCertPath := filepath.Join(pkiHomeDir, "server-cert.pem")
	serverKeyPath := filepath.Join(pkiHomeDir, "server-key.pem")
	testOrg := "localhost"
	hosts := []string{"localhost"}

	err := GenerateCert(hosts, serverCertPath, serverKeyPath, caCertPath, caKeyPath, testOrg, bits)
	if err == nil {
		_, err = os.Stat(serverCertPath)
		_, err = os.Stat(serverKeyPath)
	}

	return nil
}

func (p *PKI) Configure() (tls.Config, error) {

	config := tls.Config{}

	// TODO: Autoimport/discover CA & CACerts from MIRROR_HOME/pki/cas?
	return config, nil
}

func getTLSConfig(caCert, cert, key []byte, allowInsecure bool) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = allowInsecure
	certPool := x509.NewCertPool()

	certPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = certPool
	keypair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{keypair}
	if allowInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	return &tlsConfig, nil
}
