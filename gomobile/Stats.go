package gomobile

import "github.com/AndreasSko/go-jwlm/model"

// DatabaseStats represents the rough number of entries
// within a Database{} by defining it as the length
// of the slices.
type DatabaseStats struct {
	BlockRange int
	Bookmark   int
	Location   int
	Note       int
	Tag        int
	TagMap     int
	UserMark   int
}

// Stats generates a DatabaseStats for the given mergeSide
func (dbw *DatabaseWrapper) Stats(side string) *DatabaseStats {
	var db *model.Database

	switch side {
	case "leftSide":
		db = dbw.left
	case "rightSide":
		db = dbw.right
	case "mergeSide":
		db = dbw.merged
	default:
		db = nil
	}

	if db == nil {
		return &DatabaseStats{}
	}

	return &DatabaseStats{
		BlockRange: len(db.BlockRange),
		Bookmark:   len(db.Bookmark),
		Location:   len(db.Location),
		Note:       len(db.Note),
		Tag:        len(db.Tag),
		TagMap:     len(db.TagMap),
		UserMark:   len(db.UserMark),
	}
}
