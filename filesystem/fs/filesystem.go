package fs

import (
	"errors"
	"github.com/mefellows/mirror/filesystem"
	"io/ioutil"
	"os"
	"time"
)

// Basic File System implementation using OOTB Golang constructs
type StdFileSystem struct {
	tree filesystem.FileTree // Returns a FileTree structure of Files representing the FileSystem hierarchy
}

func (fs StdFileSystem) Dir(dir string) ([]filesystem.File, error) {
	readFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		files := make([]filesystem.File, len(readFiles))

		for i, file := range readFiles {
			files[i] = file
		}
		return files, nil
	} else {
		return nil, err
	}
}
func (fs StdFileSystem) Read(f filesystem.File) ([]byte, error) {
	return ioutil.ReadFile(f.Name())
}

func (fs StdFileSystem) Delete(file filesystem.File) error {
	return errors.New("Function not yet implemented")
}

func (fs StdFileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(file.Name(), data, perm)
}

func (fs StdFileSystem) FileTree() filesystem.FileTree {
	return nil
}

type StdFile struct {
	StdName    string      // base name of the file
	StdPath    string      // base name of the file
	StdSize    int64       // length in bytes for regular files; system-dependent for others
	StdModTime time.Time   // modification time
	StdMode    os.FileMode // File details including perms
	StdIsDir   bool        // abbreviation for Mode().IsDir()
}

func (f StdFile) Name() string {
	return f.StdName
}

func (f StdFile) Size() int64 {
	return f.StdSize
}
func (f StdFile) ModTime() time.Time {
	return f.StdModTime
}
func (f StdFile) IsDir() bool {
	return f.StdIsDir
}
func (f StdFile) Mode() os.FileMode {
	return f.StdMode
}
func (f StdFile) Sys() interface{} {
	return nil
}
