package filesystem

import (
//	"fmt"
)

// A FileTree of Files represented as a linked FileTree data-structure
type FileTree interface {
	ParentNode() FileTree
	ChildNodes() []FileTree
	File() File
}

// Diff two FileTrees - the result of the diff will be two Trees:
//
// 1. `update` containing only the updates/new additions required to be made to the target Tree
// 2. `delete` containing only the deletions required on the target Tree
//
// If a client was then to perform the corresponding updates and deletions on the target Tree
// it would then be identical in structure to the src Tree.
//
// It is up to the client to decide how to act on this information
//
// The default diffing algorithm uses modification time (`ModificationTimeFileComparator`) to determine whether or not the file is different.
// Different comparison strategies may be employed by the client (for instance, S3 may prefer to use hashes).
//
func FileTreeDiff(src FileTree, target FileTree, comparator FileComparator) (update FileTree, delete FileTree, err error) {
	// TODO: Implement a tree diff algorithm
	return nil, nil, nil
}

// Recursively walk a FileTree and run a self-type function on each node.
//
// Navigates the tree in a top left to bottom right fashion
func FileTreeWalk(tree FileTree, treeFunc func(tree FileTree) (FileTree, error)) error {
	if len(tree.ChildNodes()) > 0 {
		for _, node := range tree.ChildNodes() {
			treeFunc(node)
			FileTreeWalk(node, treeFunc)
		}
	}
	return nil
}