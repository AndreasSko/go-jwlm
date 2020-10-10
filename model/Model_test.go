package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
