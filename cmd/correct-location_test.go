package cmd

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_correctLocation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	left := &model.Database{}
	err := left.ImportJWLBackup("testdata/left.jwlibrary")
	assert.NoError(t, err)

	right := &model.Database{}
	err = right.ImportJWLBackup("testdata/right.jwlibrary")
	assert.NoError(t, err)

	uniqueKeyToLocationPreMerge := getUniqueKeyToLocationMap([]*model.Database{left, right})

	merged := model.Database{}

	// Locations
	mergedLocations, locationIDChanges, err := merger.MergeLocations(left.Location, right.Location)
	assert.NoError(t, err)
	merged.Location = mergedLocations
	merger.UpdateLRIDs(left.Bookmark, right.Bookmark, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.Bookmark, right.Bookmark, "PublicationLocationID", locationIDChanges)
	merger.UpdateLRIDs(left.InputField, right.InputField, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.Note, right.Note, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.TagMap, right.TagMap, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.UserMark, right.UserMark, "LocationID", locationIDChanges)

	// Bookmarks
	bookmarksConflictSolution := map[string]merger.MergeSolution{}
	for {
		mergedBookmarks, _, err := merger.MergeBookmarks(left.Bookmark, right.Bookmark, bookmarksConflictSolution)
		if err == nil {
			merged.Bookmark = mergedBookmarks
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			newSolutions := chooseRandomSide(err.Conflicts, &merged)
			addToSolutions(bookmarksConflictSolution, newSolutions)
		default:
			log.Fatal(err)
		}
	}

	// InputFields
	inputFieldConflictSolution := map[string]merger.MergeSolution{}
	for {
		mergedInputFields, _, err := merger.MergeInputFields(left.InputField, right.InputField, inputFieldConflictSolution)
		if err == nil {
			merged.InputField = mergedInputFields
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			newSolutions := chooseRandomSide(err.Conflicts, &merged)
			addToSolutions(inputFieldConflictSolution, newSolutions)
		default:
			log.Fatal(err)
		}
	}

	// Tags
	var tagsConflictSolution map[string]merger.MergeSolution
	for {
		mergedTags, tagIDChanges, err := merger.MergeTags(left.Tag, right.Tag, tagsConflictSolution)
		if err == nil {
			merged.Tag = mergedTags
			merger.UpdateLRIDs(left.TagMap, right.TagMap, "TagID", tagIDChanges)
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			tagsConflictSolution = chooseRandomSide(err.Conflicts, nil) // TODO
		default:
			log.Fatal(err)
		}
	}

	// Markings
	UMBRConflictSolution := map[string]merger.MergeSolution{}
	for {
		mergedUserMarks, mergedBlockRanges, userMarkIDChanges, err := merger.MergeUserMarkAndBlockRange(left.UserMark, left.BlockRange, right.UserMark, right.BlockRange, UMBRConflictSolution)
		if err == nil {
			merged.UserMark = mergedUserMarks
			merged.BlockRange = mergedBlockRanges
			merger.UpdateLRIDs(left.Note, right.Note, "UserMarkID", userMarkIDChanges)
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			newSolutions := chooseRandomSide(err.Conflicts, &merged)
			addToSolutions(UMBRConflictSolution, newSolutions)
		default:
			log.Fatal(err)
		}
	}

	// Notes
	notesConflictSolution := map[string]merger.MergeSolution{}
	for {
		mergedNotes, notesIDChanges, err := merger.MergeNotes(left.Note, right.Note, notesConflictSolution)
		if err == nil {
			merged.Note = mergedNotes
			merger.UpdateLRIDs(left.TagMap, right.TagMap, "NoteID", notesIDChanges)
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			newSolutions := chooseRandomSide(err.Conflicts, &merged)
			addToSolutions(notesConflictSolution, newSolutions)
		default:
			log.Fatal(err)
		}
	}

	// TagMaps
	var tagMapsConflictSolution map[string]merger.MergeSolution
	for {
		mergedTagMaps, _, err := merger.MergeTagMaps(left.TagMap, right.TagMap, tagMapsConflictSolution)
		if err == nil {
			merged.TagMap = mergedTagMaps
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			tagMapsConflictSolution = chooseRandomSide(err.Conflicts, nil)
		default:
			log.Fatal(err)
		}
	}

	uniqueKeyToLocationPostMerge := getUniqueKeyToLocationMap([]*model.Database{&merged})

	problems := 0
	for key, postLoc := range uniqueKeyToLocationPostMerge {
		if preLoc, ok := uniqueKeyToLocationPreMerge[key]; ok {
			if strings.HasPrefix(key, "*model.Bookmark") {
				continue
			}
			if preLoc.Equals(postLoc) {
				continue
			}
			fmt.Printf("Found difference for %s: Was '%v' is now '%v'\n", key, preLoc, postLoc)
			problems++
		}
	}
	assert.Equal(t, 0, problems)
}

func getUniqueKeyToLocationMap(dbs []*model.Database) map[string]*model.Location {
	uniqueKeyToLocation := map[string]*model.Location{}
	for _, db := range dbs {
		dbFields := reflect.ValueOf(db).Elem()
		for i := 0; i < dbFields.NumField(); i++ {
			field := dbFields.Field(i)
			if !field.CanInterface() {
				continue
			}

			tp := field.Kind()
			switch tp {
			case reflect.Slice:
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					if elem.IsNil() {
						continue
					}
					mdl := elem.Interface().(model.Model)
					uniqueKey := fmt.Sprintf("%T:%s", mdl, mdl.UniqueKey())

					mdlReflect := reflect.ValueOf(mdl).Elem()
					locIDRefl := mdlReflect.FieldByName("LocationID")
					if !locIDRefl.IsValid() {
						continue
					}
					var locID int64
					switch locIDRefl.Interface().(type) {
					case int:
						locID = locIDRefl.Int()
					case sql.NullInt32:
						val := locIDRefl.Field(0)
						locID = val.Int()
					default:
						panic("wrong type")
					}
					location := db.Location[locID]
					if location == nil {
						continue
					}
					uniqueKeyToLocation[uniqueKey] = location
				}
			default:
			}
		}
	}

	return uniqueKeyToLocation
}

func chooseRandomSide(conflicts map[string]merger.MergeConflict, mergedDB *model.Database) map[string]merger.MergeSolution {
	result := make(map[string]merger.MergeSolution, len(conflicts))
	for key, conflict := range conflicts {

		if rand.Intn(2) == 0 {
			result[key] = merger.MergeSolution{
				Side:      merger.LeftSide,
				Solution:  conflict.Left,
				Discarded: conflict.Right,
			}
		} else {
			result[key] = merger.MergeSolution{
				Side:      merger.RightSide,
				Solution:  conflict.Right,
				Discarded: conflict.Left,
			}
		}
	}

	return result
}
