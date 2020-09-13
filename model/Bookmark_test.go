package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookmark_SetID(t *testing.T) {
	m1 := &Bookmark{BookmarkID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.BookmarkID)

	m2 := Bookmark{BookmarkID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.BookmarkID)
}

func TestBookmark_UniqueKey(t *testing.T) {
	m1 := &Bookmark{
		BookmarkID:            1,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}

	assert.Equal(t, "3_4", m1.UniqueKey())
}

func TestBookmark_Equals(t *testing.T) {
	m1 := &Bookmark{
		BookmarkID:            1,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}
	m1_1 := &Bookmark{
		BookmarkID:            1000,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}
	m2 := &Bookmark{
		BookmarkID:            1000,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  5,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}

	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
}
