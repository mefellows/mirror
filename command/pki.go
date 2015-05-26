package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/pki"
	"strings"
	"time"
)

type PkiCommand struct {
	Meta             Meta
	caHost           string
	outputCA         bool
	importClientCert string
	importClientKey  string
	outputClientCert bool
	outputClientKey  bool
	importCA         string
	generateCert     bool
	configure        bool
	removePKI        bool
}

func (c *PkiCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("pki", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.caHost, "caHost", "localhost", "Specify the CAs custom hostname")
	cmdFlags.StringVar(&c.importCA, "importCA", "", "Path to CA Cert to import")
	cmdFlags.StringVar(&c.importClientCert, "importClientCert", "", "Path of client certificate to import and set as the default")
	cmdFlags.StringVar(&c.importClientKey, "importClientKey", "", "Path of client key to import and set as the default")
	cmdFlags.BoolVar(&c.configure, "configure", false, "Configures a default PKI infrastructure. Warning: This will clear any existing PKI files")
	cmdFlags.BoolVar(&c.removePKI, "removePKI", false, "Remove existing PKI keys and certs. Warning: This will require trust to be setup amongst other mirror nodes")
	cmdFlags.BoolVar(&c.outputCA, "outputCA", false, "Output the CA Certificate of this mirror node")
	cmdFlags.BoolVar(&c.outputClientCert, "outputClientCert", false, "Output the Client Certificate")
	cmdFlags.BoolVar(&c.outputClientKey, "outputClientKey", false, "Output the Client Key")
	cmdFlags.BoolVar(&c.generateCert, "generateCert", false, "Generate a custom cert from this mirror nodes' CA")

	pki, err := pki.New()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Unable to setup public key infrastructure: %s", err.Error()))
		return 1
	}

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		c.Meta.Ui.Output("invalid...")
		return 1
	}

	if c.configure {
		c.Meta.Ui.Output(fmt.Sprintf("Setting up PKI for %s...", c.caHost))
		pki.RemovePKI()
		err := pki.SetupPKI(c.caHost)
		if err != nil {
			c.Meta.Ui.Error(err.Error())
		}
		c.Meta.Ui.Output("PKI setup complete.")
	}

	if c.importCA != "" {
		c.Meta.Ui.Output(fmt.Sprintf("Importing CA from %s", c.importCA))
		timestamp := time.Now().Unix()
		err := pki.ImportCA(fmt.Sprintf("%d", timestamp), c.importCA)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Failed to import CA: %s", err.Error()))
		} else {
			c.Meta.Ui.Info("CA successfully imported")
		}
	}

	if c.importClientCert != "" && c.importClientKey != "" {
		err := pki.ImportClientCertAndKey(c.importClientCert, c.importClientKey)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Failed to import client keys: %s", err.Error()))
		} else {
			c.Meta.Ui.Info("Client keys successfully imported")
		}
	}
	if c.outputCA {
		cert, _ := pki.OutputCACert()
		c.Meta.Ui.Output(cert)
	}

	if c.outputClientCert {
		cert, _ := pki.OutputClientCert()
		c.Meta.Ui.Output(cert)
	}

	if c.outputClientKey {
		cert, _ := pki.OutputClientKey()
		c.Meta.Ui.Output(cert)
	}

	if c.removePKI {
		c.Meta.Ui.Output("Removing existing PKI")
		err := pki.RemovePKI()
		if err != nil {
			c.Meta.Ui.Error(err.Error())
		}
		c.Meta.Ui.Output("PKI removal complete.")
	}

	if c.generateCert {
		c.Meta.Ui.Output("Generating a new client cert")
		err := pki.GenerateClientCertificate([]string{"localhost"})
		if err != nil {
			c.Meta.Ui.Error(err.Error())
		}
		c.Meta.Ui.Output("Cert generation complete")
	}

	return 0
}

func (c *PkiCommand) Help() string {
	helpText := `
Usage: mirror pki [options] 

  Sets up the PKI infrastructure for secure communication between mirror nodes.
  
Options:

  --configure                 (Re-)configure PKI infrastructure on this Mirror node. This is generally only required if something strange happens to $MIRROR_HOME.
  --caHost                    Specify a custom CA Host when generating the PKI.
  --importCA                  Trust the provided CA (often used to trust other Mirror node CAs). Requires a CA Certificate.
  --outputCA                  Output the CA Certificate for this mirror node. 
  --importClientCert          Import the current Client Certificate (.crt). Must be accompanied by --importClientKey.
  --importClientKey           Import the current Client Key (.pem) file. Must be accompanied by --importClientCert.
  --outputClientCert          Output the current Client Certificate (.crt).
  --outputClientKey           Output the current Client Key (.pem) file.
  --generateCert              Generate a client cert trusted by this Mirror nodes CA.
  --removePKI                 Removes existing PKI.
`

	return strings.TrimSpace(helpText)
}

func (c *PkiCommand) Synopsis() string {
	return "Setup the PKI infrastructure for secure communication between mirror nodes."
}
