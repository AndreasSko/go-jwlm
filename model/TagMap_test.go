package model

import (
	"database/sql"
	"encoding/json"
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

func TestTagMap_RelatedEntries(t *testing.T) {
	m1 := &TagMap{
		TagMapID:       1,
		PlaylistItemID: sql.NullInt32{1, true},
		LocationID:     sql.NullInt32{1, true},
		NoteID:         sql.NullInt32{1, true},
		TagID:          1,
		Position:       1,
	}

	assert.Empty(t, m1.RelatedEntries(nil))
	assert.Empty(t, m1.RelatedEntries(&Database{}))
}

func TestTagMap_MarshalJSON(t *testing.T) {
	m1 := &TagMap{
		TagMapID:       1,
		PlaylistItemID: sql.NullInt32{2, true},
		LocationID:     sql.NullInt32{3, true},
		NoteID:         sql.NullInt32{4, true},
		TagID:          5,
		Position:       6,
	}

	result, err := json.Marshal(m1)
	assert.NoError(t, err)
	assert.Equal(t,
		`{"Type":"TagMap","TagMapID":1,"PlaylistItemID":{"Int32":2,"Valid":true},"LocationID":{"Int32":3,"Valid":true},"NoteID":{"Int32":4,"Valid":true},"TagID":5,"Position":6}`,
		string(result))
}
