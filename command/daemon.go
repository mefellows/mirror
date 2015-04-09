package command

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem/remote"
	"github.com/mefellows/mirror/pki"
	"log"
	"net"
	"net/rpc"
	//"fmt"
	"strings"
)

type DaemonCommand struct {
	Meta     Meta
	Port     int  // Which port to listen on
	Insecure bool // Enable/Disable TLS
}

func (c *DaemonCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.IntVar(&c.Port, "port", 8123, "The http port to listen on")
	cmdFlags.BoolVar(&c.Insecure, "insecure", false, "Disable TLS connection")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Would run on port %d", c.Port))
	remoteFs := new(remote.RemoteFileSystem)
	rpc.Register(remoteFs)

	service := fmt.Sprintf(":%d", c.Port)
	pkiMgr := pki.New()
	config, err := pkiMgr.GetServerTLSConfig()
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	listener, err := tls.Listen("tcp", service, config)
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
  -insecure					  Disable SSL security on the connection
`

	return strings.TrimSpace(helpText)
}

func (c *DaemonCommand) Synopsis() string {
	return "Run the mirror daemon"
}
