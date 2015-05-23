package filesystem

import (
	"errors"
	"fmt"
	fs "github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/mirror"
	"strings"
)

// TODO: This is still StdFS Specific
// Also, we need the prefix/protocol on this for MakeFile.

// Create a File struct for a new, possibly (likely) non-existent file
// but based on the attributes of a given, real file
// This is handy for creating a target File for syncing to a destination
// Note that this expects any protocol part of the path to be removed
func MkToFile(fromBase string, toBase string, file fs.File) fs.File {

	// src:  /foo/bar/baz/bat.txt
	// dest: /lol/
	// target: /lol/bat.txt

	// src: /foo/bar/baz
	// dest: s3:///lol
	// target: s3:///lol/foo/bar/baz

	// TODO: This needs some work
	path := fmt.Sprintf("%s", strings.Replace(file.Path(), fromBase, toBase, -1))
	toFile := fs.File{
		FileName:    file.Name(),
		FilePath:    path,
		FileMode:    file.Mode(),
		FileSize:    file.Size(),
		FileModTime: file.ModTime(),
	}
	return toFile
}

func GetFileSystemFromFile(file string) (fs.FileSystem, error) {
	var filesys fs.FileSystem
	var err error

	// Default protocol is "file"
	protocol := "file"
	i := strings.Index(file, "://")
	if i > -1 {
		protocol = file[:i]
	}

	// Given a protocol, find its implementor
	if factory, ok := mirror.FileSystemFactories.Lookup(protocol); ok {

		filesys, err = factory(file)
		if err != nil {
			return nil, err
		}
	} else {
		err = errors.New(fmt.Sprintf("Unable to find a suitable File System Plugin for the protocol \"%s\"", protocol))
	}
	return filesys, err

}

// Given a file path as a string (remote/local etc.) return a specific implementation
// of its filesystem.File representation, and a handle to the underling
// fs handler
func MakeFile(file string) (fs.File, fs.FileSystem, error) {
	var f fs.File
	filesys, err := GetFileSystemFromFile(file)
	if err == nil {
		f, err = filesys.ReadFile(file)
	} else {
		err = errors.New(fmt.Sprintf("Unable to find a suitable File System Plugin for the protocol: %v", err))
	}

	return f, filesys, err
}
