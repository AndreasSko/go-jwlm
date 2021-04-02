//+build js

package wasm

import (
	"fmt"
	"os"

	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// BookmarkResolver represents a resolver that should be used for conflicting Bookmarks
var BookmarkResolver string

// MarkingResolver represents a resolver that should be used for conflicting UserMarkBlockRanges
var MarkingResolver string

// NoteResolver represents a resolver that should be used for conflicting Notes
var NoteResolver string

func Merge(leftFile []byte, rightFile []byte, mergedFilename string) []byte {
	BookmarkResolver = "chooseLeft" //chooseNewest chooseLeft chooseRight
	MarkingResolver = "chooseLeft"
	NoteResolver = "chooseLeft"

	pers := model.GetPersistence()
	tmpPath, err := pers.CreateTempStorage("jwlBackups")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error creating temp path"))
	}
	defer pers.CleanupPath(tmpPath)
	//fmt.Println("Importing left backup", leftFilename)

	leftFullFileName := tmpPath + string(os.PathSeparator) + "left.jwlibrary"
	err = pers.StoreJWLBackup(leftFullFileName, leftFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error storing left backup"))
	}
	//leftFilename := left.StoreBackupData(leftFile)
	left := model.Database{}
	err = left.ImportJWLBackup(leftFullFileName)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Importing right backup", rightFilename)

	rightFullFileName := tmpPath + string(os.PathSeparator) + "right.jwlibrary"
	err = pers.StoreJWLBackup(rightFullFileName, rightFile)
	//rightFilename := right.StoreBackupData(rightFile)
	right := model.Database{}
	err = right.ImportJWLBackup(rightFullFileName)
	if err != nil {
		log.Fatal(err)
	}

	merged := model.Database{}

	fmt.Println("🧭 Merging Locations")
	mergedLocations, locationIDChanges, err := merger.MergeLocations(left.Location, right.Location)
	merged.Location = mergedLocations
	merger.UpdateLRIDs(left.Bookmark, right.Bookmark, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.Bookmark, right.Bookmark, "PublicationLocationID", locationIDChanges)
	merger.UpdateLRIDs(left.Note, right.Note, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.TagMap, right.TagMap, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.UserMark, right.UserMark, "LocationID", locationIDChanges)
	fmt.Println("Done.")

	fmt.Println("📑 Merging Bookmarks")
	bookmarksConflictSolution := map[string]merger.MergeSolution{}
	for {
		mergedBookmarks, _, err := merger.MergeBookmarks(left.Bookmark, right.Bookmark, bookmarksConflictSolution)
		if err == nil {
			merged.Bookmark = mergedBookmarks
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			if BookmarkResolver != "" {
				var resErr error
				newSolutions, resErr := merger.AutoResolveConflicts(err.Conflicts, BookmarkResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
				addToSolutions(bookmarksConflictSolution, newSolutions)
			} /*else {
				newSolutions := handleMergeConflict(err.Conflicts, &merged, stdio)
				addToSolutions(bookmarksConflictSolution, newSolutions)
			}*/
		default:
			log.Fatal(err)
		}
	}
	fmt.Println("Done.")

	fmt.Println("🏷  Merging Tags")
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
			//tagsConflictSolution = handleMergeConflict(err.Conflicts, nil, stdio) // TODO
		default:
			log.Fatal(err)
		}
	}
	fmt.Println("Done.")

	fmt.Println("🖍  Merging Markings")
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
			if MarkingResolver != "" {
				var resErr error
				newSolutions, resErr := merger.AutoResolveConflicts(err.Conflicts, MarkingResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
				addToSolutions(UMBRConflictSolution, newSolutions)
			} /*else {
				newSolutions := handleMergeConflict(err.Conflicts, &merged, stdio)
				addToSolutions(UMBRConflictSolution, newSolutions)
			}*/
		default:
			log.Fatal(err)
		}
	}
	fmt.Println("Done.")

	fmt.Println("📝 Merging Notes")
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
			if NoteResolver != "" {
				var resErr error
				newSolutions, resErr := merger.AutoResolveConflicts(err.Conflicts, NoteResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
				addToSolutions(notesConflictSolution, newSolutions)
			} /* else {
				newSolutions := handleMergeConflict(err.Conflicts, &merged, stdio)
				addToSolutions(notesConflictSolution, newSolutions)
			}*/
		default:
			log.Fatal(err)
		}
	}
	fmt.Println("Done.")

	fmt.Println("🏷  Merging TagMaps")
	var tagMapsConflictSolution map[string]merger.MergeSolution
	for {
		mergedTagMaps, _, err := merger.MergeTagMaps(left.TagMap, right.TagMap, tagMapsConflictSolution)
		if err == nil {
			merged.TagMap = mergedTagMaps
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			//tagMapsConflictSolution = handleMergeConflict(err.Conflicts, nil, stdio)
		default:
			log.Fatal(err)
		}
	}
	fmt.Println("Done.")

	fmt.Println("🎉 Finished merging!")

	fmt.Println("Exporting merged database")
	mergedPath := tmpPath + string(os.PathSeparator) + mergedFilename
	err = merged.ExportJWLBackup(mergedPath)
	if err != nil {
		log.Fatal(err)
	}

	_, mergedData, err := pers.GetFile(mergedPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error retrieving merged file."))
	}

	return mergedData

}

// addToSolutions adds new mergeSolutions to the existing map of mergeSolutions
func addToSolutions(solutions map[string]merger.MergeSolution, new map[string]merger.MergeSolution) {
	for key, value := range new {
		solutions[key] = value
	}
}
