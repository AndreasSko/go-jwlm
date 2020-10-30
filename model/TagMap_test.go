package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagMap_SetID(t *testing.T) {
	m1 := &TagMap{TagMapID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.TagMapID)

	m2 := TagMap{TagMapID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.TagMapID)
}

func TestTagMap_UniqueKey(t *testing.T) {
	m1 := &TagMap{
		TagMapID:       1,
		PlaylistItemID: sql.NullInt32{1, true},
		LocationID:     sql.NullInt32{1, true},
		NoteID:         sql.NullInt32{1, true},
		TagID:          1,
		Position:       1,
	}
	m1_1 := &TagMap{
		TagMapID:       5,
		PlaylistItemID: sql.NullInt32{1, true},
		LocationID:     sql.NullInt32{1, true},
		NoteID:         sql.NullInt32{1, true},
		TagID:          1,
		Position:       1,
	}
	m2 := &TagMap{
		TagMapID:   1,
		LocationID: sql.NullInt32{1, true},
		NoteID:     sql.NullInt32{1, true},
		TagID:      1,
		Position:   1,
	}
	assert.Equal(t, "1_1_1_1", m1.UniqueKey())
	assert.Equal(t, m1.UniqueKey(), m1_1.UniqueKey())
	assert.Equal(t, "0_1_1_1", m2.UniqueKey())
}

func TestTagMap_Equals(t *testing.T) {
	m1 := &TagMap{
		TagMapID:       1,
		PlaylistItemID: sql.NullInt32{1, true},
		LocationID:     sql.NullInt32{1, true},
		NoteID:         sql.NullInt32{1, true},
		TagID:          1,
		Position:       1,
	}
	m1_1 := &TagMap{
		TagMapID:       5,
		PlaylistItemID: sql.NullInt32{1, true},
		LocationID:     sql.NullInt32{1, true},
		NoteID:         sql.NullInt32{1, true},
		TagID:          1,
		Position:       1,
	}
	m2 := &TagMap{
		TagMapID:   1,
		LocationID: sql.NullInt32{1, true},
		NoteID:     sql.NullInt32{2, true},
		TagID:      1,
		Position:   1,
	}

	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
}
