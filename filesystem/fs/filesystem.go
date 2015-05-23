package fs

import (
	"errors"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/mirror"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Basic File System implementation using OOTB Golang constructs
type StdFileSystem struct {
	tree filesystem.FileTree // Returns a FileTree structure of Files representing the FileSystem hierarchy
	root string
}

func init() {
	mirror.FileSystemFactories.Register(NewStdFileSystem, "file")
}

func NewStdFileSystem(url string) (filesystem.FileSystem, error) {
	// Resolve/Validate URL?
	//return StdFileSystem{}, errors.New("Not yet implemented")
	return StdFileSystem{root: url}, nil
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
func FromFileInfo(dir string, i os.FileInfo) filesystem.File {
	path := fmt.Sprintf("%s/%s", dir, i.Name())
	file := filesystem.File{
		FileName:    i.Name(),
		FilePath:    path,
		FileMode:    i.Mode(),
		FileSize:    i.Size(),
		FileModTime: i.ModTime(),
	}
	return file
}

func (fs StdFileSystem) Read(f filesystem.File) ([]byte, error) {
	return ioutil.ReadFile(f.Path())
}

func (fs StdFileSystem) ReadFile(f string) (filesystem.File, error) {
	i, err := os.Stat(f)
	parentPath := filepath.Dir(f)
	if err != nil {
		return filesystem.File{}, err
	}
	return FromFileInfo(parentPath, i), err
}

func (fs StdFileSystem) Delete(file filesystem.File) error {
	return errors.New("Function not yet implemented")
}

func (fs StdFileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) error {
	parentPath := filepath.Dir(file.Path())
	if _, err := os.Stat(parentPath); err != nil {
		dir := filesystem.File{
			FilePath: parentPath,
			FileMode: 0755,
		}
		fs.MkDir(dir)
	}
	return ioutil.WriteFile(file.Path(), data, perm)
}

func (fs StdFileSystem) MkDir(file filesystem.File) error {
	return os.MkdirAll(file.Path(), file.Mode())
}

func (fs StdFileSystem) FileTree(file filesystem.File) filesystem.FileTree {
	if !file.IsDir() {
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
