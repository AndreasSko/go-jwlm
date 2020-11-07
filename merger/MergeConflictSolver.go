package merger

import (
	"fmt"
	"reflect"
	"time"
)

// MergeConflictSolver describes a function that is able to handle mergeConflicts semi-automatic
type MergeConflictSolver func(map[string]MergeConflict) (map[string]MergeSolution, error)

// AutoResolveConflicts resolves mergeConflicts using the resolver indicated by resolverName.
func AutoResolveConflicts(conflicts map[string]MergeConflict, resolverName string) (map[string]MergeSolution, error) {
	resolver, err := parseResolver(resolverName)
	if err != nil {
		return nil, err
	}
	if resolver == nil {
		return nil, nil
	}
	return resolver(conflicts)
}

// SolveConflictByChoosingLeft solves a MergeConflict by always choosing the left side
func SolveConflictByChoosingLeft(conflicts map[string]MergeConflict) (map[string]MergeSolution, error) {
	return solveConflictByChoosingSide(conflicts, LeftSide)
}

// SolveConflictByChoosingRight solves a MergeConflict by always choosing the right side
func SolveConflictByChoosingRight(conflicts map[string]MergeConflict) (map[string]MergeSolution, error) {
	return solveConflictByChoosingSide(conflicts, RightSide)
}

// parseResolver parses the name of the resolver and returns its function.
// If the name is empty, it returns nil.
func parseResolver(name string) (MergeConflictSolver, error) {
	if name == "" {
		return nil, nil
	}

	switch name {
	case "chooseLeft":
		return SolveConflictByChoosingLeft, nil
	case "chooseRight":
		return SolveConflictByChoosingRight, nil
	case "chooseNewest":
		return SolveConflictByChoosingNewest, nil
	}

	return nil, fmt.Errorf("%s is not a valid conflict resolver. Can be 'chooseNewest', 'chooseLeft', or 'chooseRight'", name)
}

// SolveConflictByChoosingNewest solves a MergeConflict by always choosing the newest entry,
// which is detected by the `LastModified` field. It returns an error if the field does not
// exist
func SolveConflictByChoosingNewest(conflicts map[string]MergeConflict) (map[string]MergeSolution, error) {
	solution := make(map[string]MergeSolution, len(conflicts))

	for key, value := range conflicts {
		leftModified := reflect.ValueOf(value.Left).Elem().FieldByName("LastModified")
		rightModified := reflect.ValueOf(value.Right).Elem().FieldByName("LastModified")

		if !leftModified.IsValid() || !rightModified.IsValid() {
			return nil, fmt.Errorf("Not able to use SolveConflictByChoosingNewest, as %T has no 'LastModified' field", value.Left)
		}

		leftDate, err := time.Parse("2006-01-02T15:04:05-07:00", leftModified.String())
		if err != nil {
			return nil, err
		}
		rightDate, err := time.Parse("2006-01-02T15:04:05-07:00", rightModified.String())
		if err != nil {
			return nil, err
		}

		if leftDate.After(rightDate) {
			solution[key] = MergeSolution{Side: LeftSide, Solution: value.Left, Discarded: value.Right}
		} else {
			solution[key] = MergeSolution{Side: RightSide, Solution: value.Right, Discarded: value.Left}
		}
	}

	return solution, nil
}

// solveConflictByChoosingSide solves a MergeConflict by always choosing the given MergeSide
func solveConflictByChoosingSide(conflicts map[string]MergeConflict, side MergeSide) (map[string]MergeSolution, error) {
	solution := make(map[string]MergeSolution, len(conflicts))

	for key, value := range conflicts {
		if side == LeftSide {
			solution[key] = MergeSolution{Side: LeftSide, Solution: value.Left, Discarded: value.Right}
		} else {
			solution[key] = MergeSolution{Side: RightSide, Solution: value.Right, Discarded: value.Left}
		}
	}

	return solution, nil
}

// solveEqualityMergeConflict solves conflicts that arise, if the same Model entry exists
// on both sides. For other conflicts it returns a mergeConflictError asking the caller
// to handle it.
func solveEqualityMergeConflict(conflicts map[string]MergeConflict) (map[string]MergeSolution, error) {
	solution := make(map[string]MergeSolution, len(conflicts))
	unsolvableConflicts := map[string]MergeConflict{}

	for key, value := range conflicts {
		if value.Left.Equals(value.Right) {
			solution[key] = MergeSolution{Side: LeftSide, Solution: value.Left, Discarded: value.Right}
		} else {
			unsolvableConflicts[key] = value
		}
	}

	if len(unsolvableConflicts) != 0 {
		return solution, MergeConflictError{Err: "Could not solve all conflicts", Conflicts: unsolvableConflicts}
	}

	return solution, nil
}
