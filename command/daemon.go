package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem/remote"
	"log"
	"net"
	//	s3 "github.com/mefellows/mirror/filesystem/s3"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"net/rpc"
	//"fmt"
	"io/ioutil"
	"strings"
)

type DaemonCommand struct {
	Meta Meta
	Port int // Which port to listen on
}

func (c *DaemonCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.IntVar(&c.Port, "port", 8123, "The http port to listen on")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Would run on port %d", c.Port))
	remoteFs := new(remote.RemoteFileSystem)
	rpc.Register(remoteFs)

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
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

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		// TODO: Need to generate proper client certs. For now, only validate if provided
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  certPool,
	}
	config.Rand = rand.Reader
	service := fmt.Sprintf(":%d", c.Port)
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}

	return 0
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	rpc.ServeConn(conn)
	log.Println("server: conn: closed")
}

func (c *DaemonCommand) Help() string {
	helpText := `
Usage: mirror daemon [options] 

  Run a mirror daemon that can listen for remote connections for file sync operations
  
Options:

  -port                       The http(s) port to listen on
  -secure					  Enable SSL security on the connection
`

	return strings.TrimSpace(helpText)
}

func (c *DaemonCommand) Synopsis() string {
	return "Run the mirror daemon"
}
