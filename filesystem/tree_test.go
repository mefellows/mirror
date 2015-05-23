package filesystem

import (
	"errors"
	"fmt"
	"testing"
)

func Test_FileTree(t *testing.T) {
	tree := makeFileTree()
	if tree.ParentNode() != nil {
		t.Fatal("tree.ParentNode() should be nil")
	}
	if tree.ChildNodes() == nil {
		t.Fatal("tree.ChildNodes() is nil")
	}
	if tree.File().Size() == 0 {
		t.Fatal("tree.File.Size() is nil")
	}
	if tree.File().Name() == "" {
		t.Fatal("tree.File.Name() is \"\"")
	}
}

func makeFileTree() FileTree {

	tree := &MockFileTree{}
	tree.FileFile = makeFile(true, File{})
	tree.ParentNodeFileTree = nil
	pnodes := make([]FileTree, 3)

	for i := 0; i < 3; i++ {
		nodes := make([]FileTree, 3)
		node := &MockFileTree{}
		node.ParentNodeFileTree = tree
		node.FileFile = makeFile(true, tree.File())
		for j := 0; j < 3; j++ {
			treeNode := &MockFileTree{}
			treeNode.FileFile = makeFile(false, node.File())
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

func TestFileTreeToMap(t *testing.T) {
	tree1 := makeFileTree()
	fileMap, _ := FileTreeToMap(tree1)

	if len(fileMap) != 12 {
		t.Fatalf("List should be size 12 but was %d", len(fileMap))
	}
}
func TestFileTreeDiff(t *testing.T) {
	tree := &MockFileTree{}
	tree.FileFile = makeFile(true, File{})
	tree.ParentNodeFileTree = nil
	pnodes := make([]FileTree, 3)

	for i := 0; i < 3; i++ {
		nodes := make([]FileTree, 3)
		node := &MockFileTree{}
		node.ParentNodeFileTree = tree
		node.FileFile = File{FileName: fmt.Sprintf("foo/%d", i)}
		for j := 0; j < 3; j++ {
			treeNode := &MockFileTree{}
			treeNode.FileFile = File{FileName: fmt.Sprintf("foo/%d/%d", i, j)}
			nodes[j] = treeNode
		}
		node.ChildNodesArray = nodes
		pnodes[i] = node
	}
	tree.ChildNodesArray = pnodes

	tree2 := &MockFileTree{}
	tree2.FileFile = makeFile(true, File{})
	tree2.ParentNodeFileTree = nil
	pnodes2 := make([]FileTree, 3)

	for i := 0; i < 3; i++ {
		nodes := make([]FileTree, 3)
		node := &MockFileTree{}
		node.ParentNodeFileTree = tree2
		node.FileFile = File{FilePath: fmt.Sprintf("foo2/%d", i)}
		for j := 0; j < 3; j++ {
			tree2Node := &MockFileTree{}
			tree2Node.FileFile = File{FilePath: fmt.Sprintf("foo/%d/%d", i, j)}
			nodes[j] = tree2Node
		}
		node.ChildNodesArray = nodes
		pnodes2[i] = node
	}
	tree2.ChildNodesArray = pnodes2

	var exists = func(l File, r File) bool {
		if r.Name() != "" {
			return true
		}
		return false
	}
	diff, _ := FileTreeDiff(tree, tree2, exists)
	if len(diff) != 3 {
		t.Fatalf("First 3 child nodes should be different (foo/{1..3} vs foo2/{1..3}. Got %d", len(diff))
	}
	fmt.Printf("Diff tree %v\n", diff)
}
