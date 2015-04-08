package command

import (
	"flag"
	"github.com/mefellows/mirror/mirror"
	"strings"
)

type PkiCommand struct {
	Meta         Meta
	generateCA   bool
	caHost       string
	generateCert bool
	configure    bool
	removePKI    bool
}

func (c *PkiCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("pki", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.BoolVar(&c.generateCA, "generateCA", false, "Generate a custom CA for this mirror node")
	cmdFlags.StringVar(&c.caHost, "caHost", "localhost", "Specify the CAs custom hostname")
	cmdFlags.BoolVar(&c.configure, "configure", false, "Configures a default PKI infrastructure")
	cmdFlags.BoolVar(&c.removePKI, "removePKI", false, "Remove existing PKI keys and certs. Warning: This will require trust to be setup amongst other mirror nodes")
	cmdFlags.BoolVar(&c.generateCert, "generateCert", false, "Generate a custom cert from this mirror nodes' CA")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if c.configure || c.generateCA {
		c.Meta.Ui.Output("Setting up PKI...")
		err := mirror.SetupPKI(c.caHost)
		if err != nil {
			c.Meta.Ui.Error(err.Error())
		}
		c.Meta.Ui.Output("PKI setup complete.")
	}

	if c.removePKI {
		c.Meta.Ui.Output("Removing existing PKI")
		err := mirror.RemovePKI()
		if err != nil {
			c.Meta.Ui.Error(err.Error())
		}
		c.Meta.Ui.Output("PKI removal complete.")
	}

	if c.generateCert {
		c.Meta.Ui.Output("Generating a client cert")
		err := mirror.GenerateCert()
		//err := mirror.GenerateCert([]string{"localhost"})
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

  --configure                 Setup PKI infrastructure on this Mirror node.
  --trustCa                   Trust the provided CA (often used to trust other Mirror node CAs).
  --generateCert              Generate a client cert trusted by this Mirror nodes CA.
  --removePKI                 Removes existing PKI.
`

	return strings.TrimSpace(helpText)
}

func (c *PkiCommand) Synopsis() string {
	return "Setup the PKI infrastructure for secure communication between mirror nodes."
}
