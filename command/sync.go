package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	utils "github.com/mefellows/mirror/filesystem/utils"
	"strings"
)

type excludes []string

func (e *excludes) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *excludes) Set(value string) error {
	*e = append(*e, value)
	return nil
}

type SyncCommand struct {
	Meta    Meta
	Dest    string
	Src     string
	Filters []string
	Exclude excludes
}

func (c *SyncCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Src, "src", "", "The src location to copy from")
	cmdFlags.StringVar(&c.Dest, "dest", "", "The destination location to copy the contents of 'src' to.")
	cmdFlags.Var(&c.Exclude, "exclude", "Set of exclusions as POSIX regular expressions to exclude from the transfer")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Syncing contents of '%s' -> '%s'", c.Src, c.Dest))

	// Obviously, this can be optimised to buffer reads directly into a write, instead of a copy and then write
	// Possibly, pass the reader into the writer and do it that way?
	fromFile, fromFs, err := utils.MakeFile(c.Src)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error opening src file: %v", err))
		return 1
	}

	if fromFile.IsDir() {
		toFile, toFs, err := utils.MakeFile(c.Dest)
		diff, err := filesystem.FileTreeDiff(
			fromFs.FileTree(fromFile), toFs.FileTree(toFile), filesystem.ModifiedComparator)

		if err == nil {
			for _, file := range diff {
				toFile = utils.MkToFile(c.Src, c.Dest, file)

				if err == nil {
					if file.IsDir() {
						//log.Printf("Mkdir: %s -> %s\n", file.Path(), toFile.Path())
						toFs.MkDir(toFile)
					} else {
						//log.Printf("Copying file: %s -> %s\n", file.Path(), toFile.Path())
						bytes, err := fromFs.Read(file)
						//log.Printf("Read bytes: %s\n", len(bytes))
						err = toFs.Write(toFile, bytes, file.Mode())
						if err != nil {
							c.Meta.Ui.Error(fmt.Sprintf("Error copying file %s: %v", file.Path(), err))
						}
					}
				}
			}
		} else {
			c.Meta.Ui.Error(fmt.Sprintf("Error: %v\n", err))
		}
	} else {
		toFile := utils.MkToFile(c.Src, c.Dest, fromFile)
		toFs, err := utils.GetFileSystemFromFile(c.Dest)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error opening dest file: %v", err))
			return 1
		}

		bytes, err := fromFs.Read(fromFile)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error reading from source file: %s", err.Error()))
			return 1
		}
		err = toFs.Write(toFile, bytes, 0644)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error writing to remote path: %s", err.Error()))
			return 1
		}
	}

	return 0
}

func (c *SyncCommand) Help() string {
	helpText := `
Usage: mirror sync [options] 

  Copy the contents of the source directory (-src) to the destination directory (-dest) recursively.
  
Options:

  -src                       The source directory from which to copy from
  -dest                      The destination directory from which to copy to
  -whatif                    Runs the sync operation as a dry-run (similar to the -n rsync flag)
  -exclude                   A regular expression used to exclude files and directories that match. 
                             This is a special option that may be specified multiple times
`

	return strings.TrimSpace(helpText)
}

func (c *SyncCommand) Synopsis() string {
	return "Copy the contents of a source directory to a destination directory"
}
