package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem/remote"
	"log"
	"net"
	//	s3 "github.com/mefellows/mirror/filesystem/s3"
	"net/http"
	"net/rpc"
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
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

	return 0
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
