package remote

import (
	"crypto/tls"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	"github.com/mefellows/mirror/filesystem/fs"
	"github.com/mefellows/mirror/mirror"
	"github.com/mefellows/mirror/pki"
	"log"
	"net"
	"net/rpc"
	neturl "net/url"
	"os"
	"strconv"
	"strings"
)

// A Remote FileSystem is essentially a Mirror daemon
// running on a remote host that can be communicated to
// over an HTTP(s) connection. Under the hood, it delegates
// to the StdFileSystem File System but is wrapped in a Go
// RPC server.
type RemoteFileSystem struct {
	rootUrl neturl.URL
	// TODO: Embed the RPC Client in here and wrap the write
	client *rpc.Client
}

func init() {
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "ssh")
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "http")
	mirror.FileSystemFactories.Register(NewRemoteFileSystem, "mirror")
}

func NewRemoteFileSystem(url string) (filesystem.FileSystem, error) {
	// Resolve/Validate URL?
	uri, err := neturl.Parse(url)

	if err != nil {
		return nil, err
	}

	// Check for host:port part
	host := uri.Host
	p := "8123"

	if strings.Contains(host, ":") {
		host, p, err = net.SplitHostPort(host)
	}
	port, _ := strconv.Atoi(p)

	// Create RPC server
	var client *rpc.Client
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), pki.MirrorConfig.ClientTlsConfig)

	// TODO: How to terminate connection when done - we want to keep open during course of events?
	//defer conn.Close()
	if err != nil {
		log.Fatalf("client: dial: %s", err)
		conn.Close()
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	client = rpc.NewClient(conn)
	return RemoteFileSystem{rootUrl: *uri, client: client}, err
}

// Remote RPC Types
type RemoteResponse struct {
	Success bool
	Error   error
}

type WriteRequest struct {
	File filesystem.File
	Data []byte
	Perm os.FileMode
}

type WriteResponse struct {
	RemoteResponse
}

type ReadFileRequest struct {
	File string
}

type ReadFileResponse struct {
	RemoteResponse
	File filesystem.File
}

type ReadRequest struct {
	File filesystem.File
}

type ReadResponse struct {
	RemoteResponse
	Data []byte
}

type FileMapRequest struct {
	File filesystem.File
}

type FileMapResponse struct {
	RemoteResponse
	FileMap filesystem.FileMap
}

type FileTreeRequest struct {
	File filesystem.File
}

type FileTreeResponse struct {
	RemoteResponse
	FileTree *filesystem.FileTree
}

type DirRequest struct {
	File string
}

type DirResponse struct {
	RemoteResponse
	Files []filesystem.File
}

type DeleteRequest struct {
	File filesystem.File
}

type DeleteResponse struct {
	RemoteResponse
}

type MkDirRequest struct {
	File filesystem.File
}

type MkDirResponse struct {
	RemoteResponse
}

func (f *RemoteFileSystem) RemoteWrite(req *WriteRequest, res *RemoteResponse) error {
	fsys := fs.StdFileSystem{}
	res.Error = fsys.Write(req.File, req.Data, req.Perm)
	if res.Error == nil {
		res.Success = true
	}
	log.Printf("Writing to file on remote side: %s @ %s. Error? %v\n", req.File.Name(), req.File.Path())
	return res.Error
}

func (f RemoteFileSystem) Write(file filesystem.File, data []byte, perm os.FileMode) (err error) {
	// Perform remote operation
	rpcargs := &WriteRequest{file, data, perm}
	var reply RemoteResponse
	err = f.client.Call("RemoteFileSystem.RemoteWrite", rpcargs, &reply)
	return err
}

func (f RemoteFileSystem) RemoteRead(req *ReadRequest, res *ReadResponse) error {
	fsys := fs.StdFileSystem{}
	res.Data, res.Error = fsys.Read(req.File)
	return res.Error
}

func (f RemoteFileSystem) Read(file filesystem.File) ([]byte, error) {
	rpcargs := &ReadRequest{File: file}
	var reply ReadResponse
	f.client.Call("RemoteFileSystem.RemoteRead", rpcargs, &reply)

	return reply.Data, reply.Error
}

func (f RemoteFileSystem) RemoteFileMap(req *FileMapRequest, res *FileMapResponse) error {
	fsys, err := fs.NewStdFileSystem(f.rootUrl.Path)
	if err == nil {
		res.FileMap = fsys.FileMap(req.File)
	} else {
		res.Error = err
	}
	return res.Error
}

func (f RemoteFileSystem) FileMap(file filesystem.File) filesystem.FileMap {
	rpcargs := &FileMapRequest{File: file}
	var reply FileMapResponse
	f.client.Call("RemoteFileSystem.RemoteFileMap", rpcargs, &reply)
	return reply.FileMap
}

func (f RemoteFileSystem) RemoteFileTree(req *FileTreeRequest, res *FileTreeResponse) error {
	fsys := fs.StdFileSystem{}
	res.FileTree = fsys.FileTree(req.File)
	return res.Error
}

func (f RemoteFileSystem) FileTree(file filesystem.File) *filesystem.FileTree {
	rpcargs := &FileTreeRequest{File: file}
	var reply FileTreeResponse
	f.client.Call("RemoteFileSystem.RemoteFileTree", rpcargs, &reply)

	return reply.FileTree
}

func (f RemoteFileSystem) RemoteDelete(req *DeleteRequest, res *DeleteResponse) error {
	fsys := fs.StdFileSystem{}
	res.Error = fsys.Delete(req.File)
	return res.Error
}

func (f RemoteFileSystem) Delete(file filesystem.File) error {
	rpcargs := &DeleteRequest{File: file}
	var reply DeleteResponse
	f.client.Call("RemoteFileSystem.RemoteDelete", rpcargs, &reply)

	return reply.Error
}

func (f RemoteFileSystem) RemoteDir(req *DirRequest, res *DirResponse) error {
	fsys := fs.StdFileSystem{}
	res.Files, res.Error = fsys.Dir(req.File)
	return res.Error
}

func (f RemoteFileSystem) Dir(dir string) ([]filesystem.File, error) {
	rpcargs := &DirRequest{File: dir}
	var reply DirResponse
	err := f.client.Call("RemoteFileSystem.RemoteDir", rpcargs, &reply)

	return reply.Files, err
}

func (f RemoteFileSystem) RemoteReadFile(req *ReadFileRequest, res *ReadFileResponse) error {
	fsys := fs.StdFileSystem{}
	res.File, res.Error = fsys.ReadFile(req.File)
	return res.Error
}

func (f RemoteFileSystem) ReadFile(file string) (filesystem.File, error) {
	rpcargs := &ReadFileRequest{File: file}
	var reply ReadFileResponse
	err := f.client.Call("RemoteFileSystem.RemoteReadFile", rpcargs, &reply)

	return reply.File, err
}

func (f RemoteFileSystem) RemoteMkdir(req *MkDirRequest, res *MkDirResponse) error {
	fsys := fs.StdFileSystem{}
	res.Error = fsys.MkDir(req.File)
	if res.Error == nil {
		res.Success = true
	}
	return res.Error
}
func (f RemoteFileSystem) MkDir(file filesystem.File) error {
	rpcargs := &MkDirRequest{File: file}
	var reply MkDirResponse
	err := f.client.Call("RemoteFileSystem.RemoteMkDir", rpcargs, &reply)

	return err
}
