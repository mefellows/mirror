package filesystem

import (
	"testing"
	"time"
)

func TestModifiedComparator(t *testing.T) {
	comp := &ModifiedComparator{}
	oldFile := &MockFile{}
	oldFile.MockIsDir = false
	oldFile.MockSize = 1024
	oldFile.MockName = "bar"
	oldFile.MockPath = "/foo/bar"
	oldFile.MockModTime = time.Now()

	newFile := &MockFile{}
	newFile.MockIsDir = false
	newFile.MockSize = 1024
	newFile.MockName = "bar"
	newFile.MockPath = "/foo/bar"
	newFile.MockModTime = time.Now()

	res := comp.Compare(newFile, oldFile)
	if !res.IsDifferent {
		t.Fatalf("Expect files to be different, got newFile: %s and oldFile: %s", newFile.ModTime(), oldFile.ModTime())
	}
}
