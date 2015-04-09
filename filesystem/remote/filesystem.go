package remote

import (
	"fmt"
	//"github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/filesystem/fs"
	"os"
)

// A Remote FileSystem is essentially a Mirror daemon
// running on a remote host that can be communicated to
// over an HTTP(s) connection. Under the hood, it uses
// the StdFileSystem File System but is wrapped in a Go
// RPC server.
type RemoteFileSystem struct {
}

type WriteRequest struct {
	// TODO: How do we make this protocol agnostic??
	// "type not registered for interface: fs.StdFile"
	//File filesystem.File
	File fs.StdFile
	Data []byte
	Perm os.FileMode
}
type WriteResponse struct {
	Success bool
}

func (f *RemoteFileSystem) Write(req *WriteRequest, res *WriteResponse) error {
	res.Success = true
	fmt.Printf("Calling remote Write!")
	fsys := fs.StdFileSystem{}
	return fsys.Write(req.File, req.Data, req.Perm)
}
