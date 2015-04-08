package mirror

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mefellows/mirror/pki"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Check for custom Certificate Authorities

// Check for server certificates (must be present for TLS to occur safely)

// Generate a client cert?

type PKI struct {
	certificateFile string
	keyFile         string
	caFile          string
}

func New() *PKI {
	return nil
}

func RemovePKI() error {
	// TODO: Get Base Config containing Home Dir
	pkiHomeDir := "/Users/mfellows/.mirror.d/pki"
	err := os.RemoveAll(pkiHomeDir)

	return err
}

//func GenerateCert(hosts []string) (err error) {
func GenerateCert() (err error) {

	outputDir := filepath.Dir(".")
	pkiHomeDir := "/Users/mfellows/.mirror.d/pki"
	caCertPath := filepath.Join(pkiHomeDir, "ca.pem")
	caKeyPath := filepath.Join(pkiHomeDir, "key.pem")
	certPath := filepath.Join(outputDir, "cert.pem")
	keyPath := filepath.Join(outputDir, "cert-key.pem")
	testOrg := "client"
	bits := 2048
	//hosts := []string{""}
	hosts := []string{}

	err = pki.GenerateCert(hosts, certPath, keyPath, caCertPath, caKeyPath, testOrg, bits)
	if err == nil {
		_, err = os.Stat(certPath)
		_, err = os.Stat(keyPath)
	}
	return
}

func SetupPKI(caHost string) error {

	// TODO: Get Base Config containing Home Dir
	pkiHomeDir := "/Users/mfellows/.mirror.d/pki"
	os.Mkdir(pkiHomeDir, 0700)

	caCertPath := filepath.Join(pkiHomeDir, "ca.pem")
	caKeyPath := filepath.Join(pkiHomeDir, "key.pem")
	bits := 2048
	if _, err := os.Stat(caCertPath); err == nil {
		return fmt.Errorf("CA already exists. Run --delete to remove the old CA.")
	}

	if err := pki.GenerateCACertificate(caCertPath, caKeyPath, caHost, bits); err != nil {
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

	err := pki.GenerateCert(hosts, serverCertPath, serverKeyPath, caCertPath, caKeyPath, testOrg, bits)
	if err == nil {
		_, err = os.Stat(serverCertPath)
		_, err = os.Stat(serverKeyPath)
	}

	return nil
}

func (p *PKI) Configure() (tls.Config, error) {

	// TODO: Autoimport/discover CA & CACerts from MIRROR_HOME/pki/cas?

	cert, err := tls.LoadX509KeyPair(p.certificateFile, p.keyFile)
	if err != nil {
		log.Fatalf("Invalid server certificate/key: %s", err.Error())
	}
	certPool := x509.NewCertPool()

	// add custom CA to pool if provided
	if p.caFile != "" {
		log.Printf("Adding custom CA (%s) to cert pool", p.caFile)
		pemData, err := ioutil.ReadFile(p.caFile)
		if err != nil {
			log.Fatalf("server: read pem file: %s", err.Error())
		}
		if ok := certPool.AppendCertsFromPEM(pemData); !ok {
			log.Fatal("server: failed to parse pem data to pool")
		}
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		Rand:         rand.Reader,
	}

	return config, err
}
