package filesystem

import (
	"os"
	"time"
)

// Generic File System abstraction
type FileSystem interface {
	Dir(string) ([]File, error)                           // Read the contents of a directory
	Read(File) ([]byte, error)                            // Read a File
	Write(file File, data []byte, perm os.FileMode) error // Write a File
	FileTree(root File) FileTree                          // Returns a FileTree structure of Files representing the FileSystem hierarchy
	MkDir(file File) error
	Delete(file File) error // Delete a file on the FileSystem
}

// Simple File abstraction (based on os.FileInfo)
//
// All local and remote files will be represented as a File.
// It is up to the specific FileSystem implementation to uphold this
//
type File interface {
	Name() string       // Name (including extension) of file
	Path() string       // Fully qualified path to file
	Size() int64        // length in bytes for regular files; system-dependent for others
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Mode() os.FileMode  // file mode bits
	Sys() interface{}   // underlying data source (can return nil)
}
