package filesystem

type FileComparator interface {
	// Compare two files and return true if they differ
	Compare(src File, dest File) FileComparisonResult
}

type FileComparisonResult struct {
	IsDifferent   bool
	IsNonExistent bool
}

// Compares the last modified time of the File
type ModifiedComparator struct{}

func (c *ModifiedComparator) Compare(src File, dest File) FileComparisonResult {
	res := &FileComparisonResult{
		IsNonExistent: (dest == nil),
	}
	if src.ModTime().After(dest.ModTime()) {
		res.IsDifferent = true
	}
	return *res
}
