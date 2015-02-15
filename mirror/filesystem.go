package mirror

import (
	"time"
)

// Generic File System abstraction
type FileSystem interface {
	Dir(string) ([]File, error) // Read the contents of a directory
	Read(File) ([]byte, error)  // Read a File
	Write(File, []byte) error   // Write a File
	Tree() Tree                 // Returns a Tree structure of Files representing the FileSystem hierarchy
}

func TreeDiff(src Tree, target Tree) (update Tree, delete Tree, err error) {
	// TODO: Implement a tree diff algorithm
	return nil, nil, nil
}

// A Tree of Files represented as a linked Tree data-structure
type Tree interface {
	ParentNode() Tree
	ChildNodes() []Tree
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
