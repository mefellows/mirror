package remote

import (
	"crypto/tls"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/filesystem/fs"
	"github.com/mefellows/mirror/mirror"
	"github.com/mefellows/mirror/pki"
	"log"
	"net/rpc"
	"os"
)

// A Remote FileSystem is essentially a Mirror daemon
// running on a remote host that can be communicated to
// over an HTTP(s) connection. Under the hood, it delegates
// to the StdFileSystem File System but is wrapped in a Go
// RPC server.
type RemoteFileSystem struct {
	root string

	// TODO: Embed the RPC Client in here and wrap the write
	client *rpc.Client
}

func init() {
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "ssh")
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "http")
}

func NewRemoteFileSystem(url string) (filesystem.FileSystem, error) {
	// Resolve/Validate URL?

	// Create RPC server
	var client *rpc.Client
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", 8123), pki.MirrorConfig.ClientTlsConfig)
	// TODO: How to terminate connection when done - we want to keep open during course of events?
	//defer conn.Close()
	if err != nil {
		log.Fatalf("client: dial: %s", err)
		conn.Close()
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	client = rpc.NewClient(conn)
	return RemoteFileSystem{root: url, client: client}, err
}

type WriteRequest struct {
	File filesystem.File
	Data []byte
	Perm os.FileMode
}

type WriteResponse struct {
	Success bool
}

func (f *RemoteFileSystem) RemoteWrite(req *WriteRequest, res *WriteResponse) error {
	fmt.Printf("Writing to file: %s @ %s\n", req.File.Name(), req.File.Path())
	log.Printf("Writing to file: %s @ %s\n", req.File.Name(), req.File.Path())
	fsys := fs.StdFileSystem{}
	err := fsys.Write(req.File, req.Data, req.Perm)
	if err == nil {
		res.Success = true
	}
	log.Printf("Writing to file on remote side: %s @ %s. Error? %v\n", req.File.Name(), req.File.Path())
	return err
}

func (f RemoteFileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) (err error) {
	fmt.Printf("Writing to file: %s @ %s", file.Name(), file.Path())
	log.Printf("Writing to file: %s @ %s", file.Name(), file.Path())

	// Perform remote operation
	rpcargs := &WriteRequest{file, data, 0644}
	var reply WriteResponse
	err = f.client.Call("RemoteFileSystem.RemoteWrite", rpcargs, &reply)

	if reply.Success {
		fmt.Printf("Hey yo - remote success!")
	}
	return err
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
