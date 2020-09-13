package merger

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/AndreasSko/go-jwlm/model"
)

// mergeSolution indicates wheter a entry came from the left or right
// side of a to-be-merged model slice pair.
type mergeSolution struct {
	side      mergeSide
	solution  model.Model
	discarded model.Model
}

type mergeSide string

const (
	leftSide  mergeSide = "leftSide"
	rightSide mergeSide = "rightSide"
)

// MergeConflict represents two Models that collide because of the same
// UniqueKey or other similarities.
type MergeConflict struct {
	left  model.Model
	right model.Model
}

// MergeConflictError indicates that a conflict happened while trying to merge
// two slices of Model. It contains the conflicts in order for the caller to solve them.
type MergeConflictError struct {
	Err       string
	Conflicts map[string]MergeConflict
}

// mergeConflictSolver describes a function that is able to handle mergeConflicts semi-automatic
type mergeConflictSolver func(map[string]MergeConflict) (map[string]mergeSolution, error)

func (e MergeConflictError) Error() string {
	return fmt.Sprintf("There were conflicts while trying to merge: %s", e.Conflicts)
}

// merge merges a left and a right slice of structs implementing the Model interface.
// If there is a collision in the process, it returns an error asking for specification how it should handle it.
func merge(left interface{}, right interface{}, conflictSolution map[string]mergeSolution) (map[string]mergeSolution, error) {
	maxLen := 0
	if reflect.ValueOf(left).Len() > reflect.ValueOf(right).Len() {
		maxLen = reflect.ValueOf(left).Len()
	} else {
		maxLen = reflect.ValueOf(right).Len()
	}

	duplicateCheck := make(map[string]mergeSolution, reflect.ValueOf(left).Len()+reflect.ValueOf(right).Len())
	collisions := make(map[string]MergeConflict, maxLen)

	// First add all entries of the left slice
	switch reflect.TypeOf(left).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(left)
		for i := 0; i < s.Len(); i++ {
			// Make sure we don't have a nil-pointer
			if s.Index(i).IsNil() {
				continue
			}

			l := s.Index(i).Interface().(model.Model)
			duplicateCheck[l.UniqueKey()] = mergeSolution{side: leftSide, solution: l}
		}
	}

	// Try to add entries of right side, if they don't conflict with existing ones
	switch reflect.TypeOf(right).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(right)
		for i := 0; i < s.Len(); i++ {
			// Make sure we don't have a nil-pointer
			if s.Index(i).IsNil() {
				continue
			}

			r := s.Index(i).Interface().(model.Model)
			if conflict, exists := duplicateCheck[r.UniqueKey()]; exists {
				if solution, ok := conflictSolution[r.UniqueKey()]; ok {
					duplicateCheck[r.UniqueKey()] = solution
				} else {
					collisions[r.UniqueKey()] = MergeConflict{
						left:  conflict.solution,
						right: r,
					}
				}
			} else {
				duplicateCheck[r.UniqueKey()] = mergeSolution{side: rightSide, solution: r}
			}
		}
	}

	if len(collisions) != 0 {
		return duplicateCheck, MergeConflictError{
			Err:       "There were conflicts while trying to merge",
			Conflicts: collisions,
		}
	}

	return duplicateCheck, nil
}

// tryMergeWithConflictSolver is a generalized method for merging a left and a right
// slice of structs implementing the Model interface. It tries to solve possible
// conflicts using the given mergeConflictSolver and will return a mergeConflictError
// if it wasn't able to solve all conflicts on its own.
func tryMergeWithConflictSolver(left interface{}, right interface{}, conflictSolution map[string]mergeSolution, conflictSolver mergeConflictSolver) ([]model.Model, IDChanges, error) {
	var solutionMap map[string]mergeSolution
	var err error

	if conflictSolution == nil {
		conflictSolution = map[string]mergeSolution{}
	}

	// Try to merge with automatic conflic resolution until the number of conflicts
	// doesn't shrink anymore
	prevConflicts := 0
Loop:
	for {
		solutionMap, err = merge(left, right, conflictSolution)
		if err == nil {
			break
		}

		switch err := err.(type) {
		case MergeConflictError:
			// If we couldn't shrink number of conflicts in last iteration, break
			if prevConflicts == len(err.Conflicts) {
				break Loop
			}

			// merge automatic conflict solution with existing/given solution
			autoConflictSolution, _ := conflictSolver(err.Conflicts)
			for key, autoSol := range autoConflictSolution {
				if sol, exists := conflictSolution[key]; exists {
					return []model.Model{}, IDChanges{}, fmt.Errorf("One of the given conflictSolution is conflicting with the one generated automatically: given %s, automatic: %s", sol, autoSol)
				}
				conflictSolution[key] = autoSol
			}

			prevConflicts = len(err.Conflicts)
		default:
			return []model.Model{}, IDChanges{}, err
		}
	}

	if err != nil {
		return []model.Model{}, IDChanges{}, err
	}

	result, changes := prepareMergeSolution(&solutionMap)

	return result, changes, err
}

// prepareMergeSolution creates are sorted slice of the solutions given in the solutionMap
// and updates the IDs of the entries if necessary. IDChanges will track changed IDs.
func prepareMergeSolution(solutionMap *map[string]mergeSolution) ([]model.Model, IDChanges) {
	// Convert map to slice and sort it so we have a deterministic output
	solutionSlice := make([]mergeSolution, len(*solutionMap))
	i := 0
	for _, sol := range *solutionMap {
		solutionSlice[i] = sol
		i++
	}
	sortMergeSolution(&solutionSlice)

	result := make([]model.Model, len(*solutionMap)+1)
	changes := IDChanges{
		Left:  map[int]int{},
		Right: map[int]int{},
	}

	i = 1
	for _, sol := range solutionSlice {
		result[i] = sol.solution
		// Update ID if needed
		if sol.solution.ID() != i {
			if sol.side == leftSide {
				changes.Left[sol.solution.ID()] = i
			} else {
				changes.Right[sol.solution.ID()] = i
			}
			result[i].SetID(i)
		}

		// If we merged a duplicate, we also need to cope with
		// changing the ID of the other side
		if sol.discarded != nil && sol.discarded.ID() != i {
			if sol.side == leftSide {
				changes.Right[sol.discarded.ID()] = i
			} else {
				changes.Left[sol.discarded.ID()] = i
			}
		}
		i++
	}

	return result, changes
}

// solveEqualityMergeConflict solves conflicts that arise, if the same Model entry exists
// on both sides. For other conflicts it returns a mergeConflictError asking the caller
// to handle it.
func solveEqualityMergeConflict(conflicts map[string]MergeConflict) (map[string]mergeSolution, error) {
	solution := make(map[string]mergeSolution, len(conflicts))
	unsolvableConflicts := map[string]MergeConflict{}

	for key, value := range conflicts {
		if value.left.Equals(value.right) {
			solution[key] = mergeSolution{side: leftSide, solution: value.left, discarded: value.right}
		} else {
			unsolvableConflicts[key] = value
		}
	}

	if len(unsolvableConflicts) != 0 {
		return solution, MergeConflictError{Err: "Could not solve all conflicts", Conflicts: unsolvableConflicts}
	}

	return solution, nil
}

// sortMergeSolution sorts a slice of mergeSolution according to the ID of the model.
// If both IDs are equal, a solution being on the left side is considered greater
// than one on the right.
func sortMergeSolution(solution *[]mergeSolution) {
	sort.Slice(*solution, func(i, j int) bool {
		// If ID are the same, then leftSide > rightSide
		if (*solution)[i].solution.ID() == (*solution)[j].solution.ID() {
			return (*solution)[i].side == leftSide
		}
		return (*solution)[i].solution.ID() < (*solution)[j].solution.ID()
	})
}
