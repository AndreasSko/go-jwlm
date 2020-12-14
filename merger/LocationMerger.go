package merger

import (
	"fmt"

	"github.com/AndreasSko/go-jwlm/model"
)

// MergeLocations merges two slices of Location into one and returns
// the merged locations together with a IDChanges struct indicating
// if the ID of a location has changed.
func MergeLocations(left []*model.Location, right []*model.Location) ([]*model.Location, IDChanges, error) {
	// Check if one side needs to migrate the bible edition from standard to study
	nwtstyMigrations := needsNwtstyMigration(left, right)
	moveToNwtsty(nwtstyMigrations, left, right)

	result, changes, err := tryMergeWithConflictSolver(left, right, nil, solveLocationMergeConflict)

	return model.Location{}.MakeSlice(result), changes, err
}

// solveLocationMergeConflict solves a merge conflict by trying to choose the Location that has
// a Title. If both don't have one, choose the right.
func solveLocationMergeConflict(conflicts map[string]MergeConflict) (map[string]MergeSolution, error) {
	solution := make(map[string]MergeSolution, len(conflicts))

	for key, value := range conflicts {
		var leftTitle string

		switch left := value.Left.(type) {
		case *model.Location:
			leftTitle = left.Title.String
		default:
			panic(fmt.Sprintf("No other type than *model.Location is supported! Given: %T", left))
		}

		if leftTitle != "" {
			solution[key] = MergeSolution{Side: LeftSide, Solution: value.Left, Discarded: value.Right}
		} else {
			solution[key] = MergeSolution{Side: RightSide, Solution: value.Right, Discarded: value.Left}
		}
	}

	return solution, nil
}

// bibleEditionCounts stores the number of occurences of the nwt/nwtsty bible
// editions on the left and the right side.
type bibleEditionCounts struct {
	leftNwt     int
	leftNwtsty  int
	rightNwt    int
	rightNwtsty int
}

// needsNwtstyMigration checks if one of the sides have been migrated from the
// Standard Bible (nwt) to the Study Edition (nwtsty), while the other one hasn't
// yet. If so, the side with the Standard Edition has to be migrated too, so
// duplicate markings can still be detected later.
//
// Side note: JW Library did the migration by setting KeySymbols from `nwt` to
// `nwtsty`, without changing the UserMarks themselfs, so their UserMarkGUID stayed the
// same. If two backups with the same markings, but with one of them not
// migrated yet, are merged, the markings can't be detected as duplicate or
// overlapping: They techically belong to different locations, though their
// UserMarkGUID is the same, which, when exporting, results in a unique
// constraint violation.
func needsNwtstyMigration(left []*model.Location, right []*model.Location) map[int]MergeSide {
	langCounts := map[int]*bibleEditionCounts{}

	for _, side := range []MergeSide{LeftSide, RightSide} {
		var entries []*model.Location
		if side == LeftSide {
			entries = left
		} else {
			entries = right
		}

		for _, location := range entries {
			if location == nil {
				continue
			}
			if _, exists := langCounts[location.MepsLanguage]; !exists {
				langCounts[location.MepsLanguage] = &bibleEditionCounts{0, 0, 0, 0}
			}

			nwt, nwtsty := 0, 0
			if location.KeySymbol.String == "nwt" {
				nwt++
			}
			if location.KeySymbol.String == "nwtsty" {
				nwtsty++
			}

			if side == LeftSide {
				langCounts[location.MepsLanguage].leftNwt += nwt
				langCounts[location.MepsLanguage].leftNwtsty += nwtsty
			} else {
				langCounts[location.MepsLanguage].rightNwt += nwt
				langCounts[location.MepsLanguage].rightNwtsty += nwtsty
			}
		}
	}

	return decideMigration(langCounts)
}

// decideMigration decides using bibleEditionCounts per mepsLanguage, if one side
// needs to be migrated from nwt to nwtsty
func decideMigration(langCounts map[int]*bibleEditionCounts) map[int]MergeSide {
	toMigrate := map[int]MergeSide{}

	for lang, counts := range langCounts {
		leftMigrated, rightMigrated := false, false

		// Consider it migrated, if the number of the study edition is way
		// higher than the number of standard edition
		if significantlyHigher(counts.leftNwtsty, counts.leftNwt) {
			leftMigrated = true
		}
		if significantlyHigher(counts.rightNwtsty, counts.rightNwt) {
			rightMigrated = true
		}

		// If both sides use same edition, we don't need to migrate any side
		if leftMigrated && rightMigrated || !leftMigrated && !rightMigrated {
			continue
		}

		if leftMigrated && !rightMigrated {
			toMigrate[lang] = RightSide
		} else {
			toMigrate[lang] = LeftSide
		}
	}

	return toMigrate
}

// moveToNwtsty moves locations from KeySymbol `nwt` to `nwtsty` if they are
// mentioned in langs by their MepsLanguage and MergeSide.
// This may be needed if both backups were started in the
// normal edition, but only one side later migrated to the study edition.
func moveToNwtsty(langs map[int]MergeSide, left []*model.Location, right []*model.Location) {
	for _, side := range []MergeSide{LeftSide, RightSide} {
		var locations []*model.Location
		if side == LeftSide {
			locations = left
		} else {
			locations = right
		}

		for _, location := range locations {
			if location == nil {
				continue
			}
			if _, exists := langs[location.MepsLanguage]; !exists {
				continue
			}
			if side != langs[location.MepsLanguage] {
				continue
			}

			if location.KeySymbol.String == "nwt" {
				location.KeySymbol.String = "nwtsty"
			}
		}
	}
}

// significantlyHigher checks if a is significantly (10x) higher than b
func significantlyHigher(a, b int) bool {
	return (float32(a) * 0.1) > float32(b)
}
