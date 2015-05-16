package filesystem

import (
	"errors"
	"sync"
)

// A FileTree is a two-way linked Tree data-structure the represents
// a hierarchy of files and directories.
// At any point in the hierarchy one can navigate freely between nodes,
// and the ordering is preserved
type FileTree interface {
	ParentNode() FileTree
	ChildNodes() []FileTree
	File() File
}

// Convert a FileTree to an Ordered ListMap
func FileTreeToMap(tree FileTree) (map[string]File, error) {

	if tree.File() == nil {
		return nil, errors.New("Empty tree")
	}

	fileMap := map[string]File{}

	treeFunc := func(tree FileTree) (FileTree, error) {
		if _, present := fileMap[tree.File().Name()]; !present {
			fileMap[tree.File().Name()] = tree.File()
		}
		return tree, nil
	}

	err := FileTreeWalk(tree, treeFunc)

	return fileMap, err
}

// Compare two file trees given a comparison function that returns true if two files are 'identical' by
// their own definition.
//
// Best we can do here is O(n) - we need to traverse 'src' and then compare 'target'
func FileTreeDiff(src FileTree, target FileTree, comparators ...func(left File, right File) bool) (diff []File, err error) {
	// Prep our two trees into lists
	var leftMap map[string]File
	var rightMap map[string]File
	var done sync.WaitGroup
	done.Add(2)
	go func() {
		leftMap, err = FileTreeToMap(src)
		done.Done()
	}()
	go func() {
		rightMap, err = FileTreeToMap(target)
		done.Done()
	}()
	done.Wait()

	// Iterate over the src list, comparing each item to the corresponding
	// match in the target Map
	diff = make([]File, 0)
	for filename, file := range leftMap {
		rightFile := rightMap[filename]
		// All comparators need to agree they are NOT different (false)
		for _, c := range comparators {
			if !c(file, rightFile) {
				diff = append(diff, file)
				break
			}
		}
	}

	return diff, nil
}

// Recursively walk a FileTree and run a self-type function on each node.
// Walker function is able to mutate the FileTree.
//
// Navigates the tree in a top left to bottom right fashion
func FileTreeWalk(tree FileTree, treeFunc func(tree FileTree) (FileTree, error)) error {
	if len(tree.ChildNodes()) > 0 {
		for _, node := range tree.ChildNodes() {

			// Mutate the tree and return any errors
			node, err := treeFunc(node)
			if err != nil {
				return err
			}
			FileTreeWalk(node, treeFunc)
		}
	}
	return nil
}
