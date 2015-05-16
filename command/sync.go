package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	fs "github.com/mefellows/mirror/filesystem/fs"
	s3 "github.com/mefellows/mirror/filesystem/s3"
	"strings"
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

	c.Meta.Ui.Output(fmt.Sprintf("Would copy contents from '%s' to '%s'", c.Src, c.Dest))

	// Obviously, this can be optimised to buffer reads directly into a write, instead of a copy and then write
	// Possibly, pass the reader into the writer and do it that way?
	fromFile, fromFs, _ := makeFile(c.Src)
	toFile, toFs, _ := makeFile(c.Dest)
	bytes, err := fromFs.Read(fromFile)
	if err != nil {
		fmt.Printf("Error reading from source file: %s", err.Error())
		return 1
	}
	err = toFs.Write(toFile, bytes, 0644)
	if err != nil {
		fmt.Printf("Error writing to remote path: %s", err.Error())
		return 1
	}

	return 0
}

// TODO: Detect File and Filesystem type
// Register FileSystems as plugins on boot?
func makeFile(file string) (filesystem.File, filesystem.FileSystem, error) {
	//
	var filesys filesystem.FileSystem
	var f filesystem.File
	var err error
	switch {
	case strings.HasPrefix(file, "s3://"):
		filesys, err = s3.New(file)
		f = s3.S3File{
			// TODO: this should be in a factory method provided by the implementor
			S3Name: file,
		}
	default:
		filesys = fs.StdFileSystem{}
		f = fs.StdFile{
			// TODO: this should be in a factory method provided by the implementor
			StdName: file,
		}
	}

	return f, filesys, err
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
