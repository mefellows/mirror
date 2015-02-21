package filesystem

import (
	"os"
	"time"
)

// Generic File System abstraction
type FileSystem interface {
	Dir(string) ([]File, error)                           // Read the contents of a directory
	Read(File) ([]byte, error)                            // Read a File
	Write(file File, data []byte, perm os.FileInfo) error // Write a File
	FileTree() FileTree                                   // Returns a FileTree structure of Files representing the FileSystem hierarchy
}

// Simple File abstraction (based on os.FileInfo)
//
// All local and remote files will be represented as a File.
// It is up to the specific FileSystem implementation to a
//
type File interface {
	Name() string // base name of the file
	//FullName() string   // fully qualified path to the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Mode() os.FileMode  // file mode bits
	Sys() interface{}   // underlying data source (can return nil)
}
