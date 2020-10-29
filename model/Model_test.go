package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeModelCopy(t *testing.T) {
	br := &BlockRange{
		BlockRangeID: 1,
		BlockType:    1,
		Identifier:   1,
		StartToken:   sql.NullInt32{Int32: 1, Valid: true},
		EndToken:     sql.NullInt32{Int32: 2, Valid: true},
		UserMarkID:   1,
	}
	brCp := MakeModelCopy(br)
	assert.Equal(t, br, brCp)
	assert.NotSame(t, br, brCp)

	bm := &Bookmark{
		BookmarkID:            1,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}
	bmCp := MakeModelCopy(bm)
	assert.Equal(t, bm, bmCp)
	assert.NotSame(t, bm, bmCp)

	loc := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   7,
		LocationType:   8,
		Title:          sql.NullString{String: "ThisTitleShouldNotBeInUniqueKey", Valid: true},
	}
	locCp := MakeModelCopy(loc)
	assert.Equal(t, loc, locCp)
	assert.NotSame(t, loc, locCp)

	note := &Note{
		NoteID:          1,
		GUID:            "GUIDFOR1",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A Title", Valid: true},
		Content:         sql.NullString{String: "The content", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	noteCp := MakeModelCopy(note)
	assert.Equal(t, note, noteCp)
	assert.NotSame(t, note, noteCp)

	tag := &Tag{
		TagID:         1,
		TagType:       1,
		Name:          "FirstTag",
		ImageFilename: sql.NullString{},
	}
	tagCp := MakeModelCopy(tag)
	assert.Equal(t, tag, tagCp)
	assert.NotSame(t, tag, tagCp)

	tm := &TagMap{
		TagMapID:       1,
		PlaylistItemID: sql.NullInt32{1, true},
		LocationID:     sql.NullInt32{1, true},
		NoteID:         sql.NullInt32{1, true},
		TagID:          1,
		Position:       1,
	}
	tmCp := MakeModelCopy(tm)
	assert.Equal(t, tm, tmCp)
	assert.NotSame(t, tm, tmCp)

	um := &UserMark{
		UserMarkID:   1,
		ColorIndex:   1,
		LocationID:   1,
		StyleIndex:   1,
		UserMarkGUID: "FIRST",
		Version:      1,
	}
	umCp := MakeModelCopy(um)
	assert.Equal(t, um, umCp)
	assert.NotSame(t, um, umCp)

	assert.Panics(t, func() {
		umbr := &UserMarkBlockRange{}
		MakeModelCopy(umbr)
	})
}

func TestPrettyPrint(t *testing.T) {
	location := &Location{
		LocationID: 1,
	}
	assert.PanicsWithValue(t, "Given struct does not contain field notexistent", func() {
		prettyPrint(location, []string{"notexistent"})
	})

	umbr := &UserMarkBlockRange{
		UserMark: &UserMark{},
	}

	assert.PanicsWithValue(t, "Unsupported type for field UserMark", func() {
		prettyPrint(umbr, []string{"UserMark"})
	})
}

func Test_sortByUniqueKey(t *testing.T) {
	locations := []*Location{
		nil,
		{
			LocationID: 1,
			BookNumber: sql.NullInt32{4, true},
		},
		{
			LocationID: 2,
			BookNumber: sql.NullInt32{3, true},
		},
		nil,
		{
			LocationID: 4,
			BookNumber: sql.NullInt32{2, true},
		},
		{
			LocationID: 5,
			BookNumber: sql.NullInt32{1, true},
		},
	}
	expectedLocations := []*Location{
		nil,
		{
			LocationID: 1,
			BookNumber: sql.NullInt32{1, true},
		},
		{
			LocationID: 2,
			BookNumber: sql.NullInt32{2, true},
		},
		{
			LocationID: 3,
			BookNumber: sql.NullInt32{3, true},
		},
		{
			LocationID: 4,
			BookNumber: sql.NullInt32{4, true},
		},
	}
	expectedLocIDChanges := map[int]int{
		1: 4,
		2: 3,
		4: 2,
		5: 1,
	}

	notes := []*Note{
		nil,
		{
			NoteID: 1,
			GUID:   "C",
		},
		{
			NoteID: 2,
			GUID:   "B",
		},
		{
			NoteID: 3,
			GUID:   "A",
		},
	}
	expectedNotes := []*Note{
		nil,
		{
			NoteID: 1,
			GUID:   "A",
		},
		{
			NoteID: 2,
			GUID:   "B",
		},
		{
			NoteID: 3,
			GUID:   "C",
		},
	}
	expectedNoteIDChanges := map[int]int{
		1: 3,
		3: 1,
	}

	locIDChanges := sortByUniqueKey(&locations)
	assert.Equal(t, expectedLocations, locations)
	assert.Equal(t, expectedLocIDChanges, locIDChanges)

	noteIDChanges := sortByUniqueKey(&notes)
	assert.Equal(t, expectedNotes, notes)
	assert.Equal(t, expectedNoteIDChanges, noteIDChanges)
}
