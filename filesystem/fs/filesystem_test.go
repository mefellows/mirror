package fs

import (
	"github.com/mefellows/mirror/filesystem"
	mirror "github.com/mefellows/mirror/mirror"
	"os"
	"testing"
)

func TestStdFileSystem(t *testing.T) {

}

func TestFileTree(t *testing.T) {

	fs := &StdFileSystem{}
	i, _ := os.Stat("/tmp/")
	file := FromFileInfo("", i)
	tree := fs.FileTree(file)

	m, _ := filesystem.FileTreeToMap(tree)
	if !(len(m) > 0) {
		t.Fatalf("Expected map size to be greater than 0")
	}
}
func TestLookupFilesystemFactory(t *testing.T) {
	filesystems := mirror.FileSystemFactories.All()
	for _, fs := range filesystems {
		f, _ := fs("foo")
		if _, ok := f.(filesystem.FileSystem); !ok {
			t.Fatalf("must be a FileSystem")
		}

	}
	f, ok := mirror.FileSystemFactories.Lookup("file")

	if !ok {
		t.Fatalf("Expected lookup to be OK")
	}

	fs, err := f("test")

	if err != nil {
		t.Fatalf("Did not expect err: %v", err)
	}

	if fs == nil {
		t.Fatalf("Expected filesystem not to be nil")
	}

}
