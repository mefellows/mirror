package mirror

import (
	"github.com/mefellows/mirror/filesystem"
	"testing"
)

func TestLookupFactory(t *testing.T) {
	NewMockFS := func(url string) (filesystem.FileSystem, error) {
		return filesystem.MockFileSystem{}, nil
	}
	FileSystemFactories.Register(NewMockFS, "mockfs")

	filesystems := FileSystemFactories.All()
	for _, fs := range filesystems {
		f, _ := fs("foo")
		if _, ok := f.(filesystem.FileSystem); !ok {
			t.Fatalf("must be a FileSystem")
		}

	}
	f, ok := FileSystemFactories.Lookup("mockfs")

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
