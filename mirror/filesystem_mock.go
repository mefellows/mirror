package mirror

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
