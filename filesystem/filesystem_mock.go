package filesystem

import (
	"os"
	"time"
)

type MockFileSystem struct {
	ReadBytes []byte
	ReadError error

	WriteError error

	DirFiles []File
	DirError error

	FileTreeTree FileTree
}

func (fs *MockFileSystem) Dir(string) ([]File, error) {
	return fs.DirFiles, fs.DirError

}
func (fs *MockFileSystem) Read(File) ([]byte, error) {
	return fs.ReadBytes, fs.ReadError
}

func (fs *MockFileSystem) Write(File, []byte) error {
	return fs.WriteError
}

func (fs *MockFileSystem) FileTree() FileTree {
	return fs.FileTreeTree
}

type MockFile struct {
	MockName     string    // base name of the file
	MockPath     string    // base name of the file
	MockSize     int64     // length in bytes for regular files; system-dependent for others
	MockModTime  time.Time // modification time
	MockIsDir    bool      // abbreviation for Mode().IsDir()
	MockFileMode os.FileMode
	MockFileSys  interface{}
}

func (f *MockFile) Name() string {
	return f.MockName
}

func (f *MockFile) Size() int64 {
	return f.MockSize
}

func (f *MockFile) ModTime() time.Time {
	return f.MockModTime
}

func (f *MockFile) IsDir() bool {
	return f.MockIsDir
}

func (f *MockFile) Mode() os.FileMode {
	return f.MockFileMode
}

func (f *MockFile) Sys() interface{} {
	return f.MockFileSys
}
