package merger

import (
	"fmt"

	"github.com/AndreasSko/go-jwlm/model"
)

// PrepareDatabasesPreMerge bundles function calls that are necessary for preparing the
// databases before merging.
func PrepareDatabasesPreMerge(left *model.Database, right *model.Database) {
	neededMigrations := needsNwtstyMigration(left, right)
	moveToNwtsty(neededMigrations, left.Location, right.Location)

	// Remove duplicate locations
	leftLocations, leftIDChanges := cleanupDuplicateLocations(left.Location)
	rightLocations, rightIDChanges := cleanupDuplicateLocations(right.Location)
	left.Location = leftLocations
	right.Location = rightLocations
	locationIDChanges := IDChanges{
		Left:  leftIDChanges,
		Right: rightIDChanges,
	}
	UpdateLRIDs(left.Bookmark, right.Bookmark, "LocationID", locationIDChanges)
	UpdateLRIDs(left.Bookmark, right.Bookmark, "PublicationLocationID", locationIDChanges)
	UpdateLRIDs(left.InputField, right.InputField, "LocationID", locationIDChanges)
	UpdateLRIDs(left.Note, right.Note, "LocationID", locationIDChanges)
	UpdateLRIDs(left.TagMap, right.TagMap, "LocationID", locationIDChanges)
	UpdateLRIDs(left.UserMark, right.UserMark, "LocationID", locationIDChanges)
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
func needsNwtstyMigration(left *model.Database, right *model.Database) map[int]MergeSide {
	// For the conflicting markings, check if one is still in `nwt`, while the other one
	// has been migrated to `nwtsty`. If that is the case, we can simply mark one side
	// to be due for migration
	leftUMGUIDs := make(map[string]*model.UserMark, len(left.UserMark))
	for _, um := range left.UserMark {
		if um == nil {
			continue
		}
		leftUMGUIDs[um.UserMarkGUID] = um
	}

	result := map[int]MergeSide{}
	for _, rightUM := range right.UserMark {
		if rightUM == nil {
			continue
		}

		leftUM, ok := leftUMGUIDs[rightUM.UserMarkGUID]
		if !ok {
			continue
		}

		leftLocation := left.Location[leftUM.LocationID]
		rightLocation := right.Location[rightUM.LocationID]
		if leftLocation.KeySymbol.String == "nwt" && rightLocation.KeySymbol.String == "nwtsty" {
			result[leftLocation.MepsLanguage] = LeftSide
			continue
		}
		if leftLocation.KeySymbol.String == "nwtsty" && rightLocation.KeySymbol.String == "nwt" {
			result[rightLocation.MepsLanguage] = RightSide
			continue
		}
	}

	return result
}

// moveToNwtsty moves locations from KeySymbol `nwt` to `nwtsty` if they are
// mentioned in langs by their MepsLanguage and MergeSide.
// This may be needed if both backups were started in the
// Regular Edition, but only one side later migrated to the Study Edition.
func moveToNwtsty(langs map[int]MergeSide, left []*model.Location, right []*model.Location) {
	if len(langs) == 0 {
		return
	}

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
			if location.KeySymbol.String != "nwt" {
				continue
			}

			// Only "simple" locations for bible books can be easily migrated.
			// For locations within a Bible that belong to a specific document,
			// we would need to find the new DocumentID. In this case, go safe
			// and skip it
			if location.DocumentID.Valid || location.Track.Valid {
				continue
			}

			location.KeySymbol.String = "nwtsty"
		}
	}
}

// cleanupDuplicateLocations looks for duplicates within one side of locations. If it finds one, it will
// choose the location that contains a title and updates the location accordingly.
// As it only checks duplicates for one side, it will directly return IDChanges in form of map[int]int.
func cleanupDuplicateLocations(entries []*model.Location) ([]*model.Location, map[int]int) {
	if entries == nil {
		return nil, nil
	}

	result := make([]*model.Location, 1, len(entries)+1)
	duplicateCheck := make(map[string]*model.Location, len(entries))
	changes := map[int]int{}

	for _, entry := range entries {
		if entry == nil {
			continue
		}
		newID := len(result)
		duplicate, ok := duplicateCheck[entry.UniqueKey()]
		if !ok {
			duplicateCheck[entry.UniqueKey()] = entry
			if entry.LocationID != newID {
				changes[entry.LocationID] = newID
				entry.LocationID = newID
			}
			result = append(result, entry)
			continue
		}

		if entry.LocationID != duplicate.LocationID {
			changes[entry.LocationID] = duplicate.LocationID
		}

		// Only replace current location with duplicate if it has a title
		if duplicate.Title.String != "" {
			continue
		}
		entry.LocationID = duplicate.LocationID
		result[duplicate.LocationID] = entry
		duplicateCheck[entry.UniqueKey()] = entry
	}

	return result, changes
}
