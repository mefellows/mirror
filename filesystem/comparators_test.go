package filesystem

import (
	"testing"
	"time"
)

func TestModifiedComparator(t *testing.T) {
	oldFile := File{
		FileSize:    1024,
		FileName:    "bar",
		FilePath:    "/foo/bar",
		FileModTime: time.Now(),
	}

	newFile := File{
		FileSize:    1024,
		FileName:    "bar",
		FilePath:    "/foo/bar",
		FileModTime: time.Now(),
	}

	res := ModifiedComparator(newFile, oldFile)
	if res {
		t.Fatalf("Expect files to be different, got newFile: %s and oldFile: %s", newFile.ModTime(), oldFile.ModTime())
	}
}
