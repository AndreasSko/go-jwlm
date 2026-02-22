package merger

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AndreasSko/go-jwlm/model"
)

// MergeIndependentMedia tries to merge the left and right slice of IndependentMedia. If there is a
// collision, it returns an error asking for specification how it should handle it.
func MergeIndependentMedia(left []*model.IndependentMedia, right []*model.IndependentMedia) ([]*model.IndependentMedia, IDChanges, error) {
	result, changes, err := tryMergeWithConflictSolver(left, right, nil, solveIndependentMediaMergeConflict)

	return model.IndependentMedia{}.MakeSlice(result), changes, err
}

func CopyMergedIndependentMedia(ims []*model.IndependentMedia, src1, src2, dst string) error {
	for _, im := range ims {
		if im == nil {
			continue
		}

		// Try to copy the file from the first source. If it doesn't exist, try the second source.
		if err := im.CopyFile(src1, dst); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err := im.CopyFile(src2, dst); err != nil {
			return err
		}
	}

	return nil
}

// solveIndependentMediaMergeConflict solves a merge conflict for IndependentMedia by choosing the side that more likely
// contains the original file name. A "." in the filename is considered an indicator that the name is the original filename.
func solveIndependentMediaMergeConflict(conflicts map[string]MergeConflict) (map[string]MergeSolution, error) {
	solution := make(map[string]MergeSolution, len(conflicts))

	for key, value := range conflicts {
		l, ok := value.Left.(*model.IndependentMedia)
		if !ok {
			return nil, fmt.Errorf("expected *model.IndependentMedia, got %T", value.Left)
		}
		r, ok := value.Right.(*model.IndependentMedia)
		if !ok {
			return nil, fmt.Errorf("expected *model.IndependentMedia, got %T", value.Right)
		}

		if l.MimeType != r.MimeType {
			return nil, fmt.Errorf("unexpected mime type mismatch for %s and %s", l, r)
		}

		if strings.Contains(l.OriginalFilename, ".") {
			solution[key] = MergeSolution{Side: LeftSide, Solution: value.Left, Discarded: value.Right}
		} else {
			solution[key] = MergeSolution{Side: RightSide, Solution: value.Right, Discarded: value.Left}
		}
	}

	return solution, nil
}
