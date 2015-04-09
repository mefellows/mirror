package pki

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
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

func New() *PKI {
	return &PKI{Config: getDefaultConfig()}
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
	testOrg := "client"
	bits := 2048

	if len(hosts) == 0 {
		hosts = []string{}
	}
	err = GenerateCert(hosts, p.Config.clientCertPath, p.Config.clientKeyPath, p.Config.caCertPath, p.Config.caKeyPath, testOrg, bits)
	if err == nil {
		_, err = os.Stat(p.Config.clientCertPath)
		_, err = os.Stat(p.Config.clientKeyPath)
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
	bits := 2048
	if _, err := os.Stat(p.Config.caCertPath); err == nil {
		return fmt.Errorf("CA already exists. Run --delete to remove the old CA.")
	}

	if err := GenerateCACertificate(p.Config.caCertPath, p.Config.caKeyPath, caHost, bits); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(p.Config.caCertPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(p.Config.caKeyPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	testOrg := "localhost"
	hosts := []string{"localhost"}

	err := GenerateCert(hosts, p.Config.serverCertPath, p.Config.serverKeyPath, p.Config.caCertPath, p.Config.caKeyPath, testOrg, bits)
	if err == nil {
		_, err = os.Stat(p.Config.serverCertPath)
		_, err = os.Stat(p.Config.serverKeyPath)
	}

	return nil
}

func (p *PKI) Configure() (tls.Config, error) {

	config := tls.Config{}

	// TODO: Autoimport/discover CA & CACerts from MIRROR_HOME/pki/cas?
	return config, nil
}

func (p *PKI) GetClientTLSConfig() (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(p.Config.clientCertPath, p.Config.clientKeyPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	pemData, err := ioutil.ReadFile(p.Config.caCertPath)
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(pemData); !ok {
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
	certPool := x509.NewCertPool()
	pemData, err := ioutil.ReadFile(p.Config.caCertPath)
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(pemData); !ok {
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
