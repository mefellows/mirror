package mirror

import (
	"time"
)

// Generic File System abstraction
type FileSystem interface {
	Dir(string) ([]File, error) // Read the contents of a directory
	Read(File) ([]byte, error)  // Read a File
	Write(File, []byte) error   // Write a File
	FileTree() FileTree         // Returns a FileTree structure of Files representing the FileSystem hierarchy
}

// TODO: Attach these Tree functions to its own class/structure/package
func FileTreeDiff(src FileTree, target FileTree) (update FileTree, delete FileTree, err error) {
	// TODO: Implement a tree diff algorithm
	return nil, nil, nil
}

// Walk a FileTree and perform some operation
func FileTreeWalk(func(*FileTree) (*FileTree, error)) error {
	return nil
}

// A FileTree of Files represented as a linked FileTree data-structure
type FileTree interface {
	ParentNode() FileTree
	ChildNodes() []FileTree
	File() File
}

// Simple File abstraction
//
// All local and remote files will be represented as a File.
// It is up to the specific FileSystem implementation to a
//
type File interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
}
