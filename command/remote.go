package command

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem/remote"
	utils "github.com/mefellows/mirror/filesystem/utils"
	"github.com/mefellows/mirror/pki"
	"log"
	"net/rpc"
	"strings"
)

type RemoteCommand struct {
	Meta     Meta
	Dest     string
	Src      string
	Host     string
	Port     int
	Cert     string
	Key      string
	Insecure bool
	Exclude  excludes
}

func (c *RemoteCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("remote", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Src, "src", "", "The src location to copy from")
	cmdFlags.StringVar(&c.Dest, "dest", "", "The destination location to copy the contents of 'src' to.")
	cmdFlags.StringVar(&c.Host, "host", "localhost", "The destination host")
	cmdFlags.StringVar(&c.Cert, "cert", "", "The location of a client certificate to use")
	cmdFlags.StringVar(&c.Key, "key", "", "The location of a client key to use")
	cmdFlags.IntVar(&c.Port, "port", 8123, "The destination host")
	cmdFlags.BoolVar(&c.Insecure, "insecure", false, "Run operation over an insecure connection")
	cmdFlags.Var(&c.Exclude, "exclude", "Set of exclusions as POSIX regular expressions to exclude from the transfer")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	pkiMgr, err := pki.New()
	pkiMgr.Config.Insecure = c.Insecure

	if c.Cert != "" {
		pkiMgr.Config.ClientCertPath = c.Cert
	}
	if c.Key != "" {
		pkiMgr.Config.ClientKeyPath = c.Key
	}

	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Unable to setup public key infrastructure: %s", err.Error()))
		return 1
	}
	config, err := pkiMgr.GetClientTLSConfig()
	if err != nil {
		log.Fatalf("Error creating TLS Config: %s", err)
		c.Meta.Ui.Error(fmt.Sprintf("Error setting up Secure communications: %s", err.Error()))
	}

	// Connect to RPC server
	var client *rpc.Client
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), config)
	defer conn.Close()
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	client = rpc.NewClient(conn)

	// Perform remote operation
	fromFile, fromFs, err := utils.MakeFile(c.Src)
	if err != nil {
		log.Printf("Error making from file: %v", err)
	}
	toFile, _, err := utils.MakeFile(c.Src)
	if err != nil {
		log.Printf("Error making to file: %v", err)
	}
	//toFs, err := utils.GetFileSystemFromFile(c.Dest)
	//fromFs, err := utils.GetFileSystemFromFile(c.Src)
	bytes, err := fromFs.Read(fromFile)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error reading from source file: %s", err.Error()))
		return 1
	}
	rpcargs := &remote.WriteRequest{toFile, bytes, 0644}
	var reply remote.WriteResponse
	err = client.Call("RemoteFileSystem.RemoteWrite", rpcargs, &reply)

	if reply.Success {
		c.Meta.Ui.Output(fmt.Sprintf("Copied '%s' to '%s'", c.Src, c.Dest))
	} else {
		c.Meta.Ui.Error(fmt.Sprintf("Unable to copy '%s' to '%s'. Error: %s", c.Src, c.Dest, err))
		return 1
	}

	return 0
}

func (c *RemoteCommand) Help() string {
	helpText := `
Usage: mirror remote [options] 

  Copy the contents of the source directory (-src) to the destination directory (-dest) recursively.
  
Options:

  --src                       The source directory from which to copy from
  --dest                      The destination directory from which to copy to
  --host                      The remote host to sync the files/folders with. Defaults to 'localhost'
  --port                      The port on the remote host to connect to. Defaults to 8123
  --insecure          		  The file transfer should be performed over an unencrypted connection
  --cert                      The certificate (.pem) to use in secure requests
  --key                       The key (.pem) to use in secure requests
`

	return strings.TrimSpace(helpText)
}

func (c *RemoteCommand) Synopsis() string {
	return "Copy the contents of a source directory to a destination directory"
}
