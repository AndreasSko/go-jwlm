package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {
	start := time.Now()

	log.Info("Importing left backup")
	left := model.Database{}
	err := left.ImportJWLBackup(os.Args[1])
	if err != nil {
		log.Warning(err)
	}

	log.Info("Importing right backup")
	right := model.Database{}
	err = right.ImportJWLBackup(os.Args[2])
	if err != nil {
		log.Warning(err)
	}
	merged := model.Database{}

	log.Info("Merging Locations")
	mergedLocations, locationIDChanges, err := merger.MergeLocations(left.Location, right.Location)
	if err != nil {
		log.Warning(err)
	}
	merged.Location = mergedLocations
	merger.UpdateIDs(left.Bookmark, right.Bookmark, "LocationID", locationIDChanges)
	merger.UpdateIDs(left.Bookmark, right.Bookmark, "PublicationLocationID", locationIDChanges)
	merger.UpdateIDs(left.Note, right.Note, "LocationID", locationIDChanges)
	merger.UpdateIDs(left.TagMap, right.TagMap, "LocationID", locationIDChanges)

	log.Info("Merging Bookmarks")
	mergedBookmarks, _, err := merger.MergeBookmarks(left.Bookmark, right.Bookmark, nil)
	if err != nil {
		log.Warning(err)
		fmt.Println("Collisions: ")
		spew.Dump(err.(merger.MergeConflictError).Conflicts)
	}
	merged.Bookmark = mergedBookmarks

	log.Info("Merging Tags")
	mergedTags, tagIDChanges, err := merger.MergeTags(left.Tag, right.Tag, nil)
	if err != nil {
		log.Warning(err)
	}
	merged.Tag = mergedTags
	merger.UpdateIDs(left.TagMap, right.TagMap, "TagID", tagIDChanges)

	log.Info("Merging TagMaps")
	mergedTagMaps, _, err := merger.MergeTagMaps(left.TagMap, right.TagMap, nil)
	if err != nil {
		log.Warning(err)
	}
	merged.TagMap = mergedTagMaps

	log.Info("Merging UserMarks & BlockRanges")
	mergedUserMarks, mergedBlockRanges, userMarkIDChanges, err := merger.MergeUserMarkAndBlockRange(left.UserMark, left.BlockRange, right.UserMark, right.BlockRange, nil) //TODO
	if err != nil {
		log.Warning(err)
		fmt.Println("Collisions: ")
		spew.Dump(err.(merger.MergeConflictError).Conflicts)
	}
	merged.UserMark = mergedUserMarks
	merged.BlockRange = mergedBlockRanges

	merger.UpdateIDs(left.Note, right.Note, "UserMarkID", userMarkIDChanges)

	log.Info("Merging Notes")
	mergedNotes, _, err := merger.MergeNotes(left.Note, right.Note, nil)
	if err != nil {
		log.Warning(err)
		fmt.Println("Collisions: ")
		spew.Dump(err.(merger.MergeConflictError).Conflicts)
	}
	merged.Note = mergedNotes

	duration := time.Since(start)
	fmt.Printf("Ran in %s", duration)
}
