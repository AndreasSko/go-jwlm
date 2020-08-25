package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNote_SetID(t *testing.T) {
	m1 := &Note{NoteID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.NoteID)

	m2 := Note{NoteID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.NoteID)
}

func TestNote_UniqueKey(t *testing.T) {
	m := &Note{NoteID: 1, GUID: "A GUID"}
	assert.Equal(t, "A GUID", m.UniqueKey())
}

func TestNote_Equals(t *testing.T) {
	m1 := &Note{
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
	m1_1 := &Note{
		NoteID:          2,
		GUID:            "GUIDFOR1",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A Title", Valid: true},
		Content:         sql.NullString{String: "The content", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	m2 := &Note{
		NoteID:          3,
		GUID:            "GUIDFOR3",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A early Title", Valid: true},
		Content:         sql.NullString{String: "The early content", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	m2_1 := &Note{
		NoteID:          3,
		GUID:            "GUIDFOR3",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A early Title", Valid: true},
		Content:         sql.NullString{String: "The early content", Valid: true},
		LastModified:    "2020-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m2_1))
}
