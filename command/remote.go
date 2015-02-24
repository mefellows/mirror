package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem/fs"
	"github.com/mefellows/mirror/filesystem/remote"
	"io/ioutil"
	//	s3 "github.com/mefellows/mirror/filesystem/s3"
	"log"
	//	"net/http"
	"crypto/tls"
	"crypto/x509"
	"net/rpc"
	"strings"
)

type RemoteCommand struct {
	Meta    Meta
	Dest    string
	Src     string
	Host    string
	Port    int
	Exclude excludes
}

func (c *RemoteCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("remote", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Src, "src", "", "The src location to copy from")
	cmdFlags.StringVar(&c.Dest, "dest", "", "The destination location to copy the contents of 'src' to.")
	cmdFlags.StringVar(&c.Host, "host", "localhost", "The destination host")
	cmdFlags.IntVar(&c.Port, "port", 8123, "The destination host")
	cmdFlags.Var(&c.Exclude, "exclude", "Set of exclusions as POSIX regular expressions to exclude from the transfer")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// Setup trust & PKI infrastructure
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	//cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	certPool := x509.NewCertPool()
	pemData, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("server: read pem file: %s", err)
	}
	if ok := certPool.AppendCertsFromPEM(pemData); !ok {
		log.Fatal("server: failed to parse pem data to pool")
	}

	// Configure TLS

	//	config := tls.Config{
	//		Certificates:       []tls.Certificate{cert},
	//		ClientAuth:         tls.RequireAndVerifyClientCert,
	//		ClientCAs:          certPool,
	//		InsecureSkipVerify: true,
	//	}
	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
	//config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	// Not passing through certs works - clearly have invalid client certs
	config = tls.Config{InsecureSkipVerify: true}

	// Connect to RPC server
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())
	client := rpc.NewClient(conn)

	// Perform remote operation
	fromFile := fs.StdFile{StdName: c.Src}
	toFile := fs.StdFile{StdName: c.Dest}
	fromFs := fs.StdFileSystem{}
	bytes, err := fromFs.Read(fromFile)
	rpcargs := &remote.WriteRequest{toFile, bytes, 0644}
	var reply remote.WriteResponse
	err = client.Call("RemoteFileSystem.Write", rpcargs, &reply)
	if err != nil {
		log.Fatal("remoteFileSystem error:", err)
	}
	fmt.Printf("Write. to file: %s, response: %s", rpcargs.File.Name(), reply)

	c.Meta.Ui.Output(fmt.Sprintf("Would copy contents from '%s' to '%s' over a remote connection", c.Src, c.Dest))

	return 0
}

func (c *RemoteCommand) Help() string {
	helpText := `
	"flag"
Usage: mirror remote [options] 

  Copy the contents of the source directory (-src) to the destination directory (-dest) recursively.
  
Options:

  -src                       The source directory from which to copy from
  -dest                      The destination directory from which to copy to
`

	return strings.TrimSpace(helpText)
}

func (c *RemoteCommand) Synopsis() string {
	return "Copy the contents of a source directory to a destination directory"
}
