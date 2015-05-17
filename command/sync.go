package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	fs "github.com/mefellows/mirror/filesystem/fs"
	s3 "github.com/mefellows/mirror/filesystem/s3"
	"os"
	"path/filepath"
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

	c.Meta.Ui.Output(fmt.Sprintf("Would copy contents from '%s' to '%s'", c.Src, c.Dest))

	// Obviously, this can be optimised to buffer reads directly into a write, instead of a copy and then write
	// Possibly, pass the reader into the writer and do it that way?
	fromFile, fromFs, err := makeFile(c.Src)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error opening src file: %v", err))
		return 1
	}
	toFile, toFs, err := makeFile(c.Dest)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error opening dest file: %v", err))
		return 1
	}

	if fromFile.IsDir() {
		diff, err := filesystem.FileTreeDiff(fromFs.FileTree(fromFile), toFs.FileTree(toFile), filesystem.ModifiedComparator)
		if err == nil {
			for _, file := range diff {
				toFile = mkToFile(c.Src, c.Dest, file)

				if err == nil {
					if file.IsDir() {
						fmt.Printf("Mkdir: %s -> %s\n", file.Path(), toFile.Path())
						toFs.MkDir(toFile)
					} else {
						fmt.Printf("Copying file: %s -> %s\n", file.Path(), toFile.Path())
						bytes, err := fromFs.Read(file)
						fmt.Printf("Read bytes: %s\n", len(bytes))
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
	}

	return 0
}

// TODO: This is still StdFS Specific
func mkToFile(fromBase string, toBase string, file filesystem.File) filesystem.File {

	// src:  /foo/bar/baz/bat.txt
	// dest: /lol/
	// target: /lol/bat.txt

	path := fmt.Sprintf("%s", strings.Replace(file.Path(), fromBase, toBase, -1))
	toFile := fs.StdFile{
		StdName:    file.Name(),
		StdPath:    path,
		StdIsDir:   file.IsDir(),
		StdMode:    file.Mode(),
		StdSize:    file.Size(),
		StdModTime: file.ModTime(),
	}
	return toFile

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
		i, err := os.Stat(file)
		if err == nil {
			f = fs.FromFileInfo(filepath.Dir(file), i)
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
