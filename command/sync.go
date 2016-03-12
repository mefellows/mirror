package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	pki "github.com/mefellows/mirror/pki"
	sync "github.com/mefellows/mirror/sync"
)

type SyncCommand struct {
	Meta     Meta
	Dest     string
	Src      string
	Host     string
	Port     int
	Cert     string
	Key      string
	Insecure bool
	Watch    bool
	Filters  []string
	Exclude  ExcludeSlice
	Verbose  bool
}

type ExcludeSlice []regexp.Regexp

func (e *ExcludeSlice) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *ExcludeSlice) Set(value string) error {
	r, err := regexp.CompilePOSIX(value)
	if err == nil {
		log.Println("posix regex:", value)
		*e = append(*e, *r)
	} else {
		log.Fatalf(fmt.Sprintf("Error parsing exclud regex: %v. Error: %v", value, err.Error()))
	}
	return nil
}

type ExcludeFlags interface {
	String() string
	Set(string) error
}

func (c *SyncCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Src, "src", "", "The src location to copy from")
	cmdFlags.StringVar(&c.Dest, "dest", "", "The destination location to copy the contents of 'src' to.")
	cmdFlags.StringVar(&c.Host, "host", "localhost", "The destination host")
	cmdFlags.StringVar(&c.Cert, "cert", "", "The location of a client certificate to use")
	cmdFlags.StringVar(&c.Key, "key", "", "The location of a client key to use")
	cmdFlags.IntVar(&c.Port, "port", 8123, "The destination host")
	cmdFlags.BoolVar(&c.Insecure, "insecure", false, "Run operation over an insecure connection")
	cmdFlags.BoolVar(&c.Watch, "watch", false, "Watch for file updates, and continuously sync on changes from --src")
	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "Enable verbose output")
	cmdFlags.Var(&c.Exclude, "exclude", "Set of exclusions as POSIX regular expressions to exclude from the transfer")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	pkiMgr, err := pki.New()

	if c.Cert != "" {
		pkiMgr.Config.ClientCertPath = c.Cert
	}
	if c.Key != "" {
		pkiMgr.Config.ClientKeyPath = c.Key
	}
	pkiMgr.Config.Insecure = c.Insecure

	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Unable to setup public key infrastructure: %s", err.Error()))
		return 1
	}
	config, err := pkiMgr.GetClientTLSConfig()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("%v", err))
		return 1
	}
	pki.MirrorConfig.ClientTlsConfig = config

	if !c.Verbose {
		log.SetOutput(ioutil.Discard)
	}
	c.Meta.Ui.Output(fmt.Sprintf("Syncing contents of '%s' -> '%s'", c.Src, c.Dest))

	options := &sync.Options{Exclude: c.Exclude}
	err = sync.Sync(c.Src, c.Dest, options)

	if c.Watch {
		c.Meta.Ui.Output(fmt.Sprintf("Monitoring %s for changes...", c.Src))
		sync.Watch(c.Src, c.Dest, options)
	}

	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error during file sync: %v", err))
		return 1
	}

	return 0
}

func (c *SyncCommand) Help() string {
	helpText := `
Usage: mirror sync [options]

  Copy the contents of the source directory (-src) to the destination directory (-dest) recursively.

Options:

  --src                       The source directory from which to copy from
  --dest                      The destination directory from which to copy to
  --whatif                    Runs the sync operation as a dry-run (similar to the -n rsync flag)
  --host                      The remote host to sync the files/folders with. Defaults to 'localhost'
  --port                      The port on the remote host to connect to. Defaults to 8123
  --insecure          		  The file transfer should be performed over an unencrypted connection
  --cert                      The certificate (.pem) to use in secure requests
  --key                       The key (.pem) to use in secure requests
  --exclude                   A regular expression used to exclude files and directories that match. Can be specified multiple times.
                              This is a special option that may be specified multiple times
  --watch                     Watch for changes in source directory and continuously sync to dest
  --verbose                   Enable output logging
`

	return strings.TrimSpace(helpText)
}

func (c *SyncCommand) Synopsis() string {
	return "Copy the contents of a source directory to a destination directory"
}
