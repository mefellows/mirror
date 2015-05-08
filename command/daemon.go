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
	"strings"
)

type DaemonCommand struct {
	Meta     Meta
	Port     int    // Which port to listen on
	Host     string // Which network host/ip to listen on
	Insecure bool   // Enable/Disable TLS
}

func (c *DaemonCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.IntVar(&c.Port, "port", 8123, "The http port to listen on")
	cmdFlags.StringVar(&c.Host, "host", "", "The host/ip to bind to. Defaults to 0.0.0.0")
	cmdFlags.BoolVar(&c.Insecure, "insecure", false, "Disable TLS connection")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Running mirror daemon on port %d", c.Port))
	remoteFs := new(remote.RemoteFileSystem)
	rpc.Register(remoteFs)

	service := fmt.Sprintf("%s:%d", c.Host, c.Port)
	pkiMgr, err := pki.New()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Unable to setup public key infrastructure: %s", err.Error()))
		return 1
	}
	pkiMgr.Config.Insecure = c.Insecure
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Unable to setup PKI infrastructure for daemon: %s", err.Error()))
		log.Fatalf("server: listen: %s", err)
	}

	var listener net.Listener
	config, err := pkiMgr.GetServerTLSConfig()
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	listener, err = tls.Listen("tcp", service, config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	for {
		conn, err := listener.Accept()
		defer conn.Close()
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

  --port                      The http(s) port to listen on
  --insecure				  Disable SSL security on the connection
`

	return strings.TrimSpace(helpText)
}

func (c *DaemonCommand) Synopsis() string {
	return "Run the mirror daemon"
}
