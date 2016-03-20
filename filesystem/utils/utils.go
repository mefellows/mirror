package filesystem

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	fs "github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/mirror"
)

func RelativeFilePath(fromBase string, toBase string, localFilePath string) string {
	return fmt.Sprintf("%s", strings.Replace(LinuxPath(localFilePath), LinuxPath(fromBase), LinuxPath(toBase), -1))
}

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
	path := fmt.Sprintf("%s", strings.Replace(LinuxPath(file.Path()), LinuxPath(fromBase), LinuxPath(toBase), -1))
	toFile := fs.File{
		FileName:    file.Name(),
		FilePath:    path,
		FileMode:    file.Mode(),
		FileSize:    file.Size(),
		FileModTime: file.ModTime(),
	}
	return toFile
}

func ExtractURL(file string) *url.URL {
	url, err := url.Parse(file)
	if err != nil {
		return nil
	}

	if url.Scheme == "" {
		url.Scheme = "file"
	}

	return url
}

func GetFileSystemFromFile(file string) (fs.FileSystem, error) {
	var filesys fs.FileSystem
	var err error

	// Default protocol is "file"
	url := ExtractURL(file)

	if url == nil {
		return nil, errors.New(fmt.Sprintf("Unable to generate URL from: %s", file))
	}

	// Given a protocol, find its implementor
	if factory, ok := mirror.FileSystemFactories.Lookup(url.Scheme); ok {

		filesys, err = factory(file)
		if err != nil {
			return nil, err
		}
	} else {
		err = errors.New(fmt.Sprintf("Unable to find a suitable File System Plugin for the protocol \"%s\"", url.Scheme))
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
		url := ExtractURL(file)
		f, err = filesys.ReadFile(url.Path)
	} else {
		err = errors.New(fmt.Sprintf("Unable to find a suitable File System Plugin for the protocol: %v", err))
	}

	return f, filesys, err
}

// TODO: Move this into a Windows specific build function and
//       make unix a no-op
func LinuxPath(path string) string {
	// Strip drive prefix c:/ etc.
	// TODO: Need to find a way to deal with paths properly (i.e. what if multiple drives!)
	r, _ := regexp.CompilePOSIX("([a-zA-Z]:)(\\.*)")
	if r.MatchString(path) {
		path = r.ReplaceAllString(path, "$2")
	}

	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, "//", "/", -1)
	return path
}
