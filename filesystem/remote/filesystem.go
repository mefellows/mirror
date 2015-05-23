package remote

import (
	"github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/filesystem/fs"
	"github.com/mefellows/mirror/mirror"
	"os"
)

// A Remote FileSystem is essentially a Mirror daemon
// running on a remote host that can be communicated to
// over an HTTP(s) connection. Under the hood, it uses
// the StdFileSystem File System but is wrapped in a Go
// RPC server.
type RemoteFileSystem struct {
	root string
}

func init() {
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "ssh")
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "http")
}

func NewRemoteFileSystem(url string) (filesystem.FileSystem, error) {
	// Resolve/Validate URL?
	//return StdFileSystem{}, errors.New("Not yet implemented")
	return RemoteFileSystem{root: url}, nil
}

type WriteRequest struct {
	// TODO: How do we make this protocol agnostic??
	// "type not registered for interface: fs.StdFile"
	//File filesystem.File
	File filesystem.File
	Data []byte
	Perm os.FileMode
}
type WriteResponse struct {
	Success bool
}

func (f *RemoteFileSystem) RemoteWrite(req *WriteRequest, res *WriteResponse) error {
	res.Success = true
	fsys := fs.StdFileSystem{}
	return fsys.Write(req.File, req.Data, req.Perm)
}

func (f RemoteFileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) error {
	fsys := fs.StdFileSystem{}
	return fsys.Write(file, data, perm)
}

func (f RemoteFileSystem) Read(file filesystem.File) ([]byte, error) {
	fsys := fs.StdFileSystem{}
	return fsys.Read(file)
}

func (f RemoteFileSystem) FileTree(file filesystem.File) filesystem.FileTree {
	fsys := fs.StdFileSystem{}
	return fsys.FileTree(file)
}

func (f RemoteFileSystem) Delete(file filesystem.File) error {
	fsys := fs.StdFileSystem{}
	return fsys.Delete(file)
}

func (f RemoteFileSystem) Dir(dir string) ([]filesystem.File, error) {
	fsys := fs.StdFileSystem{}
	return fsys.Dir(dir)
}

func (f RemoteFileSystem) ReadFile(file string) (filesystem.File, error) {
	fsys := fs.StdFileSystem{}
	return fsys.ReadFile(file)
}
func (f RemoteFileSystem) MkDir(file filesystem.File) error {
	fsys := fs.StdFileSystem{}
	return fsys.MkDir(file)
}
