package fs

import (
	"fmt"
	"github.com/mefellows/mirror/filesystem"
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

	fmt.Printf("My Tree: %v", tree)
	m, _ := filesystem.FileTreeToMap(tree)
	if !(len(m) > 0) {
		t.Fatalf("Expected map size to be greater than 0")
	}
	//fmt.Printf("My Tree, as a list")
	//for _, file := range m {
	//	fmt.Printf("tree file: %v\n", file)
	//}
}
