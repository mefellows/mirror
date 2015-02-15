package command

import (
	"flag"
	"fmt"
)

type excludes []string

func (e *excludes) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *excludes) Set(value string) error {
	fmt.Printf("%s\n", value)
	*e = append(*e, value)
	return nil
}

type SyncCommand struct {
	Meta Meta
	Dest string
	Src  string
	//Filters []string
	Exclude excludes
}

func (c *SyncCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Src, "src", "", "The src location to copy from")
	cmdFlags.StringVar(&c.Dest, "dest", "", "The destination location to copy the contents of 'src' to.")
	//cmdFlags.StringVar(&c.Filters, "filters", "", "The destination location to copy the contents of 'src' to.")
	cmdFlags.Var(&c.Exclude, "exclude", "Set of exclusions as POSIX regular expressions to exclude from the transfer")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Would copy contents from '.' to '%s'", c.Dest))
	c.Meta.Ui.Output(fmt.Sprintf("Here are the exclusions: ", c.Exclude))
	c.Meta.Ui.Error("oh shiiit")
	c.Meta.Ui.Output("Syncing from a -> b")
	c.Meta.Ui.Info("Syncing from a -> b")
	q, _ := c.Meta.Ui.Ask("Can you please tell me your age, little girl?")
	c.Meta.Ui.Info(q)
	return 0
}

func (c *SyncCommand) Help() string {
	return "Sync's a folder from A -> B"
}

func (c *SyncCommand) Synopsis() string {
	return "Run interactive"
}
