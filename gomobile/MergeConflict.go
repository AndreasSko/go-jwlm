package gomobile

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
)

// MergeConflictError indicates that a conflict happened while merging. It
// is equivalent to merger.MergeConflictError, but does not contain the
// actual conflicts to make it compatible with Gomobile.
type MergeConflictError struct {
	Err string
}

func (e MergeConflictError) Error() string {
	return fmt.Sprintf("There were conflicts while trying to merge %s", e.Err)
}

// MergeConflictsWrapper wraps mergeConflicts and their solutions
type MergeConflictsWrapper struct {
	DBWrapper       *DatabaseWrapper
	conflicts       map[string]merger.MergeConflict
	conflictKeys    []string
	solvedConflicts int
	solutions       map[string]merger.MergeSolution
}

// MergeConflict represents two Models that collide. It is equvalent
// to merger.MergeConflict, but represents the Models as strings
// to make it compatible with Gomobile.
type MergeConflict struct {
	Left  string
	Right string
}

// modelRelatedTuple contains a model and its related entries
type modelRelatedTuple struct {
	Model   model.Model
	Related model.Related
}

// InitDBWrapper initializes the DatabaseWrapper for the MergeConflictsWrapper
// so the DB is accessible for pretty printing.
func (mcw *MergeConflictsWrapper) InitDBWrapper(dbw *DatabaseWrapper) {
	mcw.DBWrapper = dbw
}

// addConflicts adds the given conflicts to the MergeConflictWrapper
// and makes sure that the keys are added to conflictsKeys. This
// makes sure, that the order stays the same. If the conflict already
// exists, it is skipped.
func (mcw *MergeConflictsWrapper) addConflicts(conflicts map[string]merger.MergeConflict) {
	if mcw.conflicts == nil {
		mcw.conflicts = make(map[string]merger.MergeConflict, len(conflicts))
	}

	if mcw.conflictKeys == nil {
		mcw.conflictKeys = make([]string, 0, len(conflicts))
	}

	for key, value := range conflicts {
		if _, exists := mcw.conflicts[key]; !exists {
			mcw.conflicts[key] = value
			mcw.conflictKeys = append(mcw.conflictKeys, key)
		}
	}
}

// ConflictsLen returns the length of the conflicts map.
func (mcw *MergeConflictsWrapper) ConflictsLen() int {
	if mcw.conflicts == nil {
		return 0
	}

	return len(mcw.conflicts)
}

// SolutionsLen returns the length of the solutions slice
func (mcw *MergeConflictsWrapper) SolutionsLen() int {
	if mcw.solutions == nil {
		return 0
	}

	return len(mcw.solutions)
}

// GetNextConflictIndex returns the next conflict index, for which
// the conflict is not solved yet. If there are none left, it returns
// -1 indicating that all conflicts have been solved.
func (mcw *MergeConflictsWrapper) GetNextConflictIndex() int {
	if mcw.solvedConflicts >= len(mcw.conflicts) {
		return -1
	}

	return mcw.solvedConflicts
}

// GetConflict returns the conflict at index
func (mcw *MergeConflictsWrapper) GetConflict(index int) (*MergeConflict, error) {
	if mcw.conflicts == nil || mcw.conflictKeys == nil {
		return nil, errors.New("There are no conflicts")
	}

	if index >= len(mcw.conflictKeys) {
		return nil, fmt.Errorf("Conflict with index %d does not exist. Length=%d", index, len(mcw.conflictKeys))
	}
	key := mcw.conflictKeys[index]

	if mcw.DBWrapper == nil {
		mcw.DBWrapper = &DatabaseWrapper{merged: nil}
	}

	result := &MergeConflict{}

	jsn, err := json.Marshal(modelRelatedTuple{
		Model:   mcw.conflicts[key].Left,
		Related: mcw.conflicts[key].Left.RelatedEntries(mcw.DBWrapper.merged),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error while marshalling to JSON")
	}
	result.Left = string(jsn)

	jsn, err = json.Marshal(modelRelatedTuple{
		Model:   mcw.conflicts[key].Right,
		Related: mcw.conflicts[key].Right.RelatedEntries(mcw.DBWrapper.merged),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error while marshalling to JSON")
	}
	result.Right = string(jsn)

	return result, nil
}

// SolveConflict solves a mergeConflict by choosing the given side at index.
// Index must be less or equal to GetNextConflictIndex(), to ensure that
// conflicts are solved in order and none are missed.
func (mcw *MergeConflictsWrapper) SolveConflict(index int, side string) error {
	if mcw.conflicts == nil || mcw.conflictKeys == nil {
		return errors.New("There are no conflicts")
	}
	if index >= len(mcw.conflictKeys) {
		return fmt.Errorf("Conflict with index %d does not exist. Length=%d", index, len(mcw.conflictKeys))
	}
	if mcw.GetNextConflictIndex() != -1 && mcw.GetNextConflictIndex() < index {
		return fmt.Errorf("Index is higher than NextConflictIndex: %d > %d. The conflicts before have to be solved first",
			index, mcw.GetNextConflictIndex())
	}

	if mcw.solutions == nil {
		mcw.solutions = make(map[string]merger.MergeSolution, len(mcw.conflictKeys))
	}

	key := mcw.conflictKeys[index]
	switch side {
	case "leftSide":
		mcw.solutions[key] = merger.MergeSolution{
			Side:      merger.LeftSide,
			Solution:  mcw.conflicts[key].Left,
			Discarded: mcw.conflicts[key].Right,
		}
	case "rightSide":
		mcw.solutions[key] = merger.MergeSolution{
			Side:      merger.RightSide,
			Solution:  mcw.conflicts[key].Right,
			Discarded: mcw.conflicts[key].Left,
		}
	default:
		return fmt.Errorf("Side %s is not valid", side)
	}

	if mcw.GetNextConflictIndex() == index {
		mcw.solvedConflicts++
	}

	return nil
}
