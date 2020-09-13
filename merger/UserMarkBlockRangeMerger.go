package merger

import (
	"fmt"
	"sort"

	"github.com/AndreasSko/go-jwlm/model"
)

// brFrom indicates from which mergeSide a *BlockRange
// is coming from.
type brFrom struct {
	side mergeSide
	br   *model.BlockRange
}

// umbrFrom indicates from which mergeSide a []*UserMarkBlockRange
// is coming from.
type umbrFrom struct {
	side mergeSide
	umbr []*model.UserMarkBlockRange
}

// MergeUserMarkAndBlockRange joins UserMarks and BlockRanges from both sides and
// tries to merge them. Afterwards it will update the IDs of UserMark and BlockRange
// and returns them separately againg. If there is a collision, it will try to solve
// it using duplicate detection and - if that fails - returns an error asking for
// specification how it should handle it. MergeConflicts will be returned as a joined
// UserMarkBlockRange struct to make it easier representing conflicts.
// The returned IDChanges indicate if a UserMarkID has changed in the merge process.
func MergeUserMarkAndBlockRange(leftUM []*model.UserMark, leftBR []*model.BlockRange,
	rightUM []*model.UserMark, rightBR []*model.BlockRange,
	conflictSolution map[string]mergeSolution) ([]*model.UserMark, []*model.BlockRange, IDChanges, error) {
	if conflictSolution == nil {
		conflictSolution = map[string]mergeSolution{}
	}

	left := joinToUserMarkBlockRange(leftUM, leftBR)
	right := joinToUserMarkBlockRange(rightUM, rightBR)

	merged, changes, err := mergeUMBR(left, right, conflictSolution)
	um, br := splitUserMarkBlockRange(merged)
	if err == nil {
		return um, br, changes, nil
	}

	// If merge failed, try to solve conflicts using solveEqualityMergeConflict
	switch err := err.(type) {
	case MergeConflictError:
		autoConflictSolution, _ := solveEqualityMergeConflict(err.Conflicts)
		for key, autoSol := range autoConflictSolution {
			conflictSolution["_"+key] = autoSol
		}
	default:
		return nil, nil, IDChanges{}, err
	}

	merged, changes, err = mergeUMBR(left, right, conflictSolution)
	um, br = splitUserMarkBlockRange(merged)

	return um, br, changes, err
}

// mergeUMBR merges a left and a right side of *[]UserMarkBlockRange. It will check
// for overlapping (i.e. conflicting) BlockRanges and returns an mergeConflictError
// if it finds some, asking the caller for specification how it should handle it.
// IDChanges indicate if a UserMarkID has changed in the merge process.
func mergeUMBR(left []*model.UserMarkBlockRange, right []*model.UserMarkBlockRange,
	conflictSolution map[string]mergeSolution) ([]*model.UserMarkBlockRange, IDChanges, error) {
	// First, replace conflictSolution entries with the conflicting ones on the left
	// and right side, so we don't detect them again.
	changes, invertedChanges := replaceUMBRConflictsWithSolution(&left, &right, conflictSolution)

	conflicts := map[string]MergeConflict{}

	// Ingest UserMarks & BlockRanges in a Map[LocationID]map[Identifier][]*model.BlockRange
	blRanges := ingestUMBR(left, right)
	conflictsCount := 0
	// For each BlockRange slice per identifier, sort BlockRanges by StartToken
	for _, locationBlock := range blRanges {
		for _, identifierBlock := range locationBlock {
			sort.SliceStable(identifierBlock, func(i, j int) bool {
				return identifierBlock[i].br.StartToken.Int32 < identifierBlock[j].br.StartToken.Int32
			})

			// Check for overlapping intervals
			for i, br := range identifierBlock {
				for j := i + 1; j < len(identifierBlock) && br.br.EndToken.Int32 >= identifierBlock[j].br.StartToken.Int32; j++ {
					// We found a collision!

					// If collision is on the same side, then ignore it
					// (it's probably not our fault and we hope its okay...)
					if br.side == identifierBlock[j].side {
						continue
					}

					// If its one different sites, add it to conflicts & make sure
					// that entries are on correct side of mergeConflict{}
					if br.side == leftSide {
						conflicts[fmt.Sprint(conflictsCount)] = MergeConflict{left[br.br.UserMarkID], right[identifierBlock[j].br.UserMarkID]}
					} else {
						conflicts[fmt.Sprint(conflictsCount)] = MergeConflict{left[identifierBlock[j].br.UserMarkID], right[br.br.UserMarkID]}
					}
					conflictsCount++
				}
			}
		}
	}

	if len(conflicts) > 0 {
		return []*model.UserMarkBlockRange{}, IDChanges{}, MergeConflictError{Conflicts: conflicts}
	}

	// Add left and right to result & update (UserMark-)ID
	result := make([]*model.UserMarkBlockRange, len(left)+len(right)+2)
	i := 1
	for _, mergeSide := range []mergeSide{leftSide, rightSide} {
		var side []*model.UserMarkBlockRange
		if mergeSide == leftSide {
			side = left
		} else {
			side = right
		}

		for _, entry := range side {
			if entry == nil {
				continue
			}
			// Note IDChanges if necessary
			if entry.ID() != i {
				if mergeSide == leftSide {
					// Check if on the other side an ID has changed to entry.ID
					// after its entry was discareded. If so, we again need update
					// the entry, as the current ID it is pointing at will be changed too.
					if val, ok := invertedChanges.Right[entry.ID()]; ok {
						changes.Right[val] = i
					}
					changes.Left[entry.ID()] = i
				} else {
					if val, ok := invertedChanges.Left[entry.ID()]; ok {
						changes.Left[val] = i
					}
					changes.Right[entry.ID()] = i
				}
			}

			result[i] = entry
			result[i].SetID(i)
			i++
		}
	}

	return result[:i], changes, nil
}

// replaceUMBRConflictsWithSolution removes conflicting entries on the left and right
// and replaces one of them with the given solution in conflictSolution. IDChanges
// and its inverted counterpart indicate possible changes of IDs.
func replaceUMBRConflictsWithSolution(left *[]*model.UserMarkBlockRange, right *[]*model.UserMarkBlockRange, conflictSolution map[string]mergeSolution) (IDChanges, IDChanges) {
	changes := IDChanges{
		Left:  map[int]int{},
		Right: map[int]int{},
	}
	invertedChanges := IDChanges{
		Left:  map[int]int{},
		Right: map[int]int{},
	}

	for _, sol := range conflictSolution {
		var side, other *[]*model.UserMarkBlockRange
		if sol.side == leftSide {
			side = left
			other = right
			changes.Right[sol.discarded.ID()] = sol.solution.ID()
			invertedChanges.Right[sol.solution.ID()] = sol.discarded.ID()
		} else {
			side = right
			other = left
			changes.Left[sol.discarded.ID()] = sol.solution.ID()
			invertedChanges.Left[sol.solution.ID()] = sol.discarded.ID()
		}

		(*side)[sol.solution.ID()] = (sol.solution).(*model.UserMarkBlockRange)
		(*other)[sol.discarded.ID()] = nil
	}

	return changes, invertedChanges
}

// ingestUMBR ingest UserMarks & BlockRanges in a Map[LocationID]map[Identifier][]*model.BlockRange
func ingestUMBR(left []*model.UserMarkBlockRange, right []*model.UserMarkBlockRange) map[int]map[int][]brFrom {
	result := make(map[int]map[int][]brFrom, estimateLocationCount(left, right))
	for _, side := range []*umbrFrom{{leftSide, left}, {rightSide, right}} {
		for _, umbr := range side.umbr {
			if umbr == nil {
				continue
			}
			if _, ok := result[umbr.UserMark.LocationID]; !ok {
				result[umbr.UserMark.LocationID] = map[int][]brFrom{}
			}
			for _, br := range umbr.BlockRanges {
				if _, ok := result[umbr.UserMark.LocationID][br.Identifier]; !ok {
					result[umbr.UserMark.LocationID][br.Identifier] = []brFrom{}
				}
				result[umbr.UserMark.LocationID][br.Identifier] = append(result[umbr.UserMark.LocationID][br.Identifier], brFrom{side.side, br})
			}
		}
	}

	return result
}

// joinToUserMarkBlockRange joins entries of UserMark and BlockRange together creating
// a slice of UserMarkBlockRange for which the index corresponds to the UserMarkID.
// It expects the IDs of UserMark and BlockRange to correspond to their
// index within in um and br!
func joinToUserMarkBlockRange(um []*model.UserMark, br []*model.BlockRange) []*model.UserMarkBlockRange {
	result := make([]*model.UserMarkBlockRange, len(um))
	for i, entry := range um {
		if entry == nil || *entry == (model.UserMark{}) {
			continue
		}
		result[i] = &model.UserMarkBlockRange{UserMark: entry, BlockRanges: []*model.BlockRange{}}
	}

	for _, entry := range br {
		if entry == nil || entry.UserMarkID >= len(result) || result[entry.UserMarkID] == nil {
			continue
		}
		result[entry.UserMarkID].BlockRanges = append(result[entry.UserMarkID].BlockRanges, entry)
	}

	return result
}

// splitUserMarkBlockRange splits a UserMarkBlockRange into separate UserMark and BlockRange slices.
// It expects the UserMark IDs of BlockRanges within a UserMarkBlockRange to not be correct yet, so it will
// update them according to their UserMark. It will also update the IDs of all entries to be in sync
// with the index of the slices.
func splitUserMarkBlockRange(umbrs []*model.UserMarkBlockRange) ([]*model.UserMark, []*model.BlockRange) {
	if len(umbrs) == 0 {
		return []*model.UserMark{}, []*model.BlockRange{}
	}

	userMarks := make([]*model.UserMark, len(umbrs))
	blockRanges := make([]*model.BlockRange, 1, len(umbrs))
	blockRanges[0] = nil

	brIndex := 1
	for i, umbr := range umbrs {
		if umbr == nil {
			continue
		}

		userMarks[i] = umbr.UserMark
		for _, br := range umbr.BlockRanges {
			if br == nil || *br == (model.BlockRange{}) {
				continue
			}
			br.UserMarkID = umbr.UserMark.UserMarkID
			br.SetID(brIndex)
			blockRanges = append(blockRanges, br)
			brIndex++
		}
	}

	return userMarks, blockRanges
}

// estimateLocationCount tries to estimate the number of locations of the left
// and right side combined by looking at the highest LocationID.
func estimateLocationCount(left []*model.UserMarkBlockRange, right []*model.UserMarkBlockRange) int {
	maxLocationID := 0
	for _, side := range [][]*model.UserMarkBlockRange{left, right} {
		for i := len(side) - 1; i >= 0; i-- {
			if side[i] == nil || side[i].UserMark.LocationID == 0 {
				continue
			}
			maxLocationID += side[i].UserMark.LocationID
			break
		}
	}
	return maxLocationID
}
