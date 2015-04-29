package pki

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/mefellows/mirror/mirror"
	"io/ioutil"
	"os"
	"path"
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
	err := os.RemoveAll(path.Dir(p.Config.caCertPath))
	if err != nil {
		return err
	}

	// Client certificates
	err = os.RemoveAll(path.Dir(p.Config.clientCertPath))
	if err != nil {
		return err
	}

	// Server certificates
	err = os.RemoveAll(path.Dir(p.Config.serverKeyPath))
	if err != nil {
		return err
	}

	return err
}

func (p *PKI) GenerateCert(hosts []string) (err error) {
	organisation := "client"
	bits := 2048

	if len(hosts) == 0 {
		hosts = []string{}
	}
	err = GenerateCert(hosts, p.Config.clientCertPath, p.Config.clientKeyPath, p.Config.caCertPath, p.Config.caKeyPath, organisation, bits)
	if err == nil {
		_, err = os.Stat(p.Config.clientCertPath)
		_, err = os.Stat(p.Config.clientKeyPath)
	}
	return
}

// Validate all components of the PKI infrastructure are properly configured
func (p *PKI) CheckSetup() error {
	var err error

	err = errors.New("Not yet implemented")

	// Check directories
	if _, err := os.Stat(p.Config.caCertPath); err == nil {
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

	bits := 2048
	if _, err := os.Stat(p.Config.caCertPath); err == nil {
		return fmt.Errorf("CA already exists. Run --delete to remove the old CA.")
	}

	os.MkdirAll(path.Dir(p.Config.caCertPath), 0700)
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

	os.MkdirAll(path.Dir(p.Config.serverCertPath), 0700)
	err := GenerateCert(hosts, p.Config.serverCertPath, p.Config.serverKeyPath, p.Config.caCertPath, p.Config.caKeyPath, organisation, bits)
	if err == nil {
		_, err = os.Stat(p.Config.serverCertPath)
		_, err = os.Stat(p.Config.serverKeyPath)
	}

	// Setup Client side...
	p.GenerateCert([]string{"localhost"})

	return nil
}

func outputFileContents(file string) string {
	f, err := ioutil.ReadFile(file)
	if err == nil {
		return string(f)
	}
	return ""

}
func (p *PKI) OutputClientKey() string {
	return outputFileContents(p.Config.clientKeyPath)
}

func (p *PKI) OutputClientCert() string {
	return outputFileContents(p.Config.clientCertPath)
}

func (p *PKI) OutputCAKey() string {
	return outputFileContents(p.Config.caKeyPath)
}
func (p *PKI) OutputCACert() string {
	return outputFileContents(p.Config.caCertPath)
}

func (p *PKI) GetClientTLSConfig() (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(p.Config.clientCertPath, p.Config.clientKeyPath)
	if err != nil {
		return nil, err
	}

	certPool, err := p.discoverCAs()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
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
