// General notes on this package:
//   1. Needs all non-PKI specific stuff removed
//   2. Needs command related logging removed
//   3. Ideally this package can be extracted from this project and made more generic
package pki

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mefellows/mirror/mirror"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// General Mirror Pubic Key Infrastructure functions

type PKI struct {
	Config *Config
}

type Config struct {
	clientKeyPath  string
	clientCertPath string
	serverKeyPath  string
	serverCertPath string
	caKeyPath      string
	caCertPath     string
	Insecure       bool
}

func NewWithConfig(config *Config) (*PKI, error) {
	pki := &PKI{Config: config}
	if err := pki.SetupPKI("localhost"); err != nil {
		return nil, err
	}
	return pki, nil
}

func New() (*PKI, error) {
	pki := &PKI{Config: getDefaultConfig()}
	if err := pki.SetupPKI("localhost"); err != nil {
		return nil, err
	}
	return pki, nil
}

func getDefaultConfig() *Config {
	caHomeDir := mirror.GetCADir()
	certDir := mirror.GetCertDir()
	caCertPath := filepath.Join(caHomeDir, "ca.pem")
	caKeyPath := filepath.Join(caHomeDir, "key.pem")
	certPath := filepath.Join(certDir, "cert.pem")
	keyPath := filepath.Join(certDir, "cert-key.pem")
	serverCertPath := filepath.Join(certDir, "server-cert.pem")
	serverKeyPath := filepath.Join(certDir, "server-key.pem")

	return &Config{
		clientKeyPath:  keyPath,
		clientCertPath: certPath,
		caCertPath:     caCertPath,
		caKeyPath:      caKeyPath,
		serverCertPath: serverCertPath,
		serverKeyPath:  serverKeyPath,
	}
}

func (p *PKI) RemovePKI() error {
	// Root CA + Certificates
	err := os.RemoveAll(filepath.Dir(p.Config.caCertPath))
	if err != nil {
		return err
	}

	// Client certificates
	err = os.RemoveAll(filepath.Dir(p.Config.clientCertPath))
	if err != nil {
		return err
	}

	// Server certificates
	err = os.RemoveAll(filepath.Dir(p.Config.serverKeyPath))
	if err != nil {
		return err
	}

	return err
}

func (p *PKI) GenerateClientCertificate(hosts []string) (err error) {
	organisation := "client"
	bits := 2048

	if len(hosts) == 0 {
		hosts = []string{}
	}
	err = GenerateCertificate(hosts, p.Config.clientCertPath, p.Config.clientKeyPath, p.Config.caCertPath, p.Config.caKeyPath, organisation, bits)
	if err == nil {
		_, err = os.Stat(p.Config.clientCertPath)
		_, err = os.Stat(p.Config.clientKeyPath)
	}
	return
}

// Validate all components of the PKI infrastructure are properly configured
func (p *PKI) CheckSetup() error {
	var err error

	// Check directories
	if _, err = os.Stat(p.Config.caCertPath); err == nil {
		return nil
	}

	// Check CA

	// Check server cert

	// Check client certs against CA (from conf + user?)

	// Check permissions?

	return err
}

// Sets up the PKI infrastructure for client / server communications
// This involves creating directories, CAs, and client/server certs
func (p *PKI) SetupPKI(caHost string) error {
	if p.CheckSetup() == nil {
		return nil
	}
	log.Printf("Setting up PKI...")

	bits := 2048
	if _, err := os.Stat(p.Config.caCertPath); err == nil {
		return fmt.Errorf("CA already exists. Run --delete to remove the old CA.")
	}

	os.MkdirAll(filepath.Dir(p.Config.caCertPath), 0700)
	if err := GenerateCACertificate(p.Config.caCertPath, p.Config.caKeyPath, caHost, bits); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(p.Config.caCertPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(p.Config.caKeyPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	organisation := "localhost"
	hosts := []string{"localhost"}

	os.MkdirAll(filepath.Dir(p.Config.serverCertPath), 0700)
	err := GenerateCertificate(hosts, p.Config.serverCertPath, p.Config.serverKeyPath, p.Config.caCertPath, p.Config.caKeyPath, organisation, bits)
	if err == nil {
		_, err = os.Stat(p.Config.serverCertPath)
		_, err = os.Stat(p.Config.serverKeyPath)
	}

	// Setup Client side...
	p.GenerateClientCertificate([]string{"localhost"})

	return nil
}

func (p *PKI) OutputClientKey() (string, error) {
	return mirror.OutputFileContents(p.Config.clientKeyPath)
}

func (p *PKI) OutputClientCert() (string, error) {
	return mirror.OutputFileContents(p.Config.clientCertPath)
}

func (p *PKI) OutputCAKey() (string, error) {
	return mirror.OutputFileContents(p.Config.caKeyPath)
}
func (p *PKI) OutputCACert() (string, error) {
	return mirror.OutputFileContents(p.Config.caCertPath)
}

func (p *PKI) GetClientTLSConfig() (*tls.Config, error) {

	var certificates []tls.Certificate
	if !p.Config.Insecure {
		cert, err := tls.LoadX509KeyPair(p.Config.clientCertPath, p.Config.clientKeyPath)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, cert)
	}

	certPool, err := p.discoverCAs()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates:       certificates,
		RootCAs:            certPool,
		InsecureSkipVerify: p.Config.Insecure,
	}

	return config, err
}

func (p *PKI) GetServerTLSConfig() (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(p.Config.serverCertPath, p.Config.serverKeyPath)
	if err != nil {
		return nil, err
	}

	certPool, err := p.discoverCAs()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		Rand:         rand.Reader,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	if p.Config.Insecure {
		fmt.Printf("Setting server auth to insecure")
		config.ClientAuth = tls.NoClientCert
	}

	return config, err
}

func (p *PKI) discoverCAs() (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	// Default Root CA
	caPaths := []string{p.Config.caCertPath}
	var err error

	for _, cert := range caPaths {
		pemData, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, err
		}
		if ok := certPool.AppendCertsFromPEM(pemData); !ok {
			return nil, err
		}
	}

	return certPool, err
}
