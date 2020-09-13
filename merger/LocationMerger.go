package merger

import (
	"fmt"

	"github.com/AndreasSko/go-jwlm/model"
)

// MergeLocations merges two slices of Location into one and returns
// the merged locations together with a IDChanges struct indicating
// if the ID of a location has changed.
func MergeLocations(left []*model.Location, right []*model.Location) ([]*model.Location, IDChanges, error) {
	result, changes, err := tryMergeWithConflictSolver(left, right, nil, solveLocationMergeConflict)

	return model.Location{}.MakeSlice(result), changes, err
}

// solveLocationMergeConflict solves a merge conflict by trying to choose the Location that has
// a Title. If both don't have one, choose the right.
func solveLocationMergeConflict(conflicts map[string]mergeConflict) (map[string]mergeSolution, error) {
	solution := make(map[string]mergeSolution, len(conflicts))

	for key, value := range conflicts {
		var leftTitle string

		switch left := value.left.(type) {
		case *model.Location:
			leftTitle = left.Title.String
		default:
			panic(fmt.Sprintf("No other type than *model.Location is supported! Given: %T", left))
		}

		if leftTitle != "" {
			solution[key] = mergeSolution{side: leftSide, solution: value.left, discarded: value.right}
		} else {
			solution[key] = mergeSolution{side: rightSide, solution: value.right, discarded: value.left}
		}
	}

	return solution, nil
}
