package filesystem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}
func Test_MockFileSystem(t *testing.T) {
	d1 := []byte("hello\ngo\n")
	ioutil.WriteFile("/tmp/dat1", d1, 0644)

	mock := &MockFileSystem{}
	mockFile := makeFile(false)
	d1 = []byte("hello\ngo\n")
	path := fmt.Sprintf("%s%s-%d", os.TempDir(), "testmockfilesystem-", time.Now().UnixNano())
	ioutil.WriteFile(path, d1, 0644)
	f, _ := os.Open(path)
	defer f.Close()
	r := bufio.NewReader(f)

	mock.ReadBytes = make([]byte, 9)
	r.Read(mock.ReadBytes)
	bytesRead, _ := mock.Read(mockFile)

	if bytes.Compare(bytesRead, d1) != 0 {
		t.Fatalf("File read should read the same set of bytes back. Expected %s, got %s", d1, bytesRead)
	}
	mock.DirError = errors.New("Directory doesn't exist")
	mock.Dir("foo")
	mock.FileTree()
	mock.Write(nil, make([]byte, 0))

}

func Test_MockFileTree(t *testing.T) {
	tree := makeFileTree()
	if tree.ParentNode() != nil {
		t.Fatal("tree.ParentNode() should be nil")
	}
	if tree.ChildNodes() == nil {
		t.Fatal("tree.ChildNodes() is nil")
	}
	if tree.File() == nil {
		t.Fatal("tree.File() is nil")
	}
	if tree.File().Size() == 0 {
		t.Fatal("tree.File.Size() is nil")
	}
	if tree.File().Name() == "" {
		t.Fatal("tree.File.Name() is \"\"")
	}
	//if tree.File().FullName() == "" {
	//	t.Fatal("tree.File.FullName() is \"\"")
	//}
}

func makeFileTree() FileTree {

	tree := &MockFileTree{}
	tree.FileFile = makeFile(true)
	tree.ParentNodeFileTree = nil
	pnodes := make([]FileTree, 3)

	for i := 0; i < 3; i++ {
		nodes := make([]FileTree, 3)
		node := &MockFileTree{}
		node.ParentNodeFileTree = tree
		node.FileFile = makeFile(true)
		for j := 0; j < 3; j++ {
			treeNode := &MockFileTree{}
			treeNode.FileFile = makeFile(false)
			nodes[j] = treeNode
		}
		node.ChildNodesArray = nodes
		pnodes[i] = node
	}
	tree.ChildNodesArray = pnodes
	return tree
}

func TestFileTreeWalk(t *testing.T) {
	tree := makeFileTree()
	count := 0
	dirCount := 0
	treeFunc := func(tree FileTree) (FileTree, error) {
		if !tree.File().IsDir() {
			count++
		} else {
			dirCount++
		}
		return tree, nil
	}
	FileTreeWalk(tree, treeFunc)
	if count != 9 {
		t.Fatalf("Expected to iterate over exactly 9 files, got %d", count)
	}
	if dirCount != 3 {
		t.Fatalf("Expected to iterate over exactly 4 directories, got %d", dirCount)
	}
}

func TestFileTreeWalk_Error(t *testing.T) {
	tree := makeFileTree()
	count := 0
	dirCount := 0
	treeFunc := func(tree FileTree) (FileTree, error) {
		if tree.File().IsDir() {
			dirCount++
		} else {
			count++
		}
		return tree, errors.New("This is expected")
	}
	err := FileTreeWalk(tree, treeFunc)
	if err == nil {
		t.Fatal("Expected err")
	}
	if count != 0 {
		t.Fatalf("Expected to iterate over exactly 0 files, got %d", count)
	}
	if dirCount != 1 {
		t.Fatalf("Expected to iterate over exactly 1 directory, got %d", dirCount)
	}

}

func makeFile(isDir bool) File {
	f := &MockFile{}
	f.MockIsDir = isDir
	f.MockSize = rand.Int63()
	f.MockName = makeRandomWord()
	f.MockPath = fmt.Sprintf("%s/%s", makeRandomWord(), makeRandomWord())
	f.MockModTime = time.Now()
	return f
}

func makeRandomWord() string {
	words := []string{"foo", "bar", "bat", "baz", "crab", "cat", "parsnip", "apple", "futon"}
	return words[rand.Intn(len(words))]
}
