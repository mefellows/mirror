package fs

import (
	"errors"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Basic File System implementation using OOTB Golang constructs
type StdFileSystem struct {
	tree filesystem.FileTree // Returns a FileTree structure of Files representing the FileSystem hierarchy
}

func (fs StdFileSystem) Dir(dir string) ([]filesystem.File, error) {
	readFiles, err := ioutil.ReadDir(fmt.Sprintf("%v/", dir))
	if err == nil {
		files := make([]filesystem.File, len(readFiles))

		for i, file := range readFiles {
			files[i] = FromFileInfo(dir, file)
		}

		return files, nil
	} else {
		return nil, err
	}
}

// Converts a FileInfo -> StdFile
func FromFileInfo(dir string, i os.FileInfo) StdFile {
	path := fmt.Sprintf("%s/%s", dir, i.Name())
	file := StdFile{
		StdName:    i.Name(),
		StdPath:    path,
		StdIsDir:   i.IsDir(),
		StdMode:    i.Mode(),
		StdSize:    i.Size(),
		StdModTime: i.ModTime(),
	}
	return file

}

func (fs StdFileSystem) Read(f filesystem.File) ([]byte, error) {
	return ioutil.ReadFile(f.Path())
}

func (fs StdFileSystem) Delete(file filesystem.File) error {
	return errors.New("Function not yet implemented")
}

func (fs StdFileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) error {
	parentPath := filepath.Dir(file.Path())
	if _, err := os.Stat(parentPath); err != nil {
		dir := &StdFile{
			StdPath: parentPath,
			StdMode: 0755,
		}
		fs.MkDir(dir)
	}
	return ioutil.WriteFile(file.Path(), data, perm)
}

func (fs StdFileSystem) MkDir(file filesystem.File) error {
	return os.MkdirAll(file.Path(), file.Mode())
}

func (fs StdFileSystem) FileTree(file filesystem.File) filesystem.FileTree {
	if file == nil || !file.IsDir() {
		return nil
	}
	tree := &filesystem.StdFileSystemTree{}
	tree.StdFile = file
	return fs.readDir(file, tree)
}

// Recursively read a directory structure and create a tree structure out of it
// TODO: fix symlinks/cyclic dependencies etc.
func (fs StdFileSystem) readDir(curFile filesystem.File, parent *filesystem.StdFileSystemTree) filesystem.FileTree {
	tree := &filesystem.StdFileSystemTree{}
	tree.StdFile = curFile
	tree.StdParentNode = parent

	// TODO: Symlink check not working...
	if curFile.IsDir() || curFile.Mode() == os.ModeSymlink {

		tree.StdChildNodes = make([]filesystem.FileTree, 0)
		dirListing, _ := fs.Dir(curFile.Path())
		if dirListing != nil && len(dirListing) > 0 {
			for _, file := range dirListing {
				tree.StdChildNodes = append(tree.StdChildNodes, fs.readDir(file, tree))
			}
		}
	}

	return tree

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

func (f StdFile) Path() string {
	return f.StdPath
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
