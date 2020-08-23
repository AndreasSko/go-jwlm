package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_SetID(t *testing.T) {
	m1 := &Tag{TagID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.TagID)

	m2 := Tag{TagID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.TagID)
}

func TestTag_UniqueKey(t *testing.T) {
	m1 := &Tag{
		TagID:         1,
		TagType:       1,
		Name:          "FirstTag",
		ImageFilename: sql.NullString{},
	}
	assert.Equal(t, "1_FirstTag", m1.UniqueKey())

	m2 := &Tag{
		TagID:         1,
		TagType:       2000000000,
		Name:          "Another Tag with spaces",
		ImageFilename: sql.NullString{},
	}
	assert.Equal(t, "2000000000_Another Tag with spaces", m2.UniqueKey())
}

func TestTag_Equals(t *testing.T) {
	m1 := &Tag{
		TagID:         1,
		TagType:       1,
		Name:          "FirstTag",
		ImageFilename: sql.NullString{},
	}
	m1_1 := &Tag{
		TagID:         100000,
		TagType:       1,
		Name:          "FirstTag",
		ImageFilename: sql.NullString{},
	}
	m2 := &Tag{
		TagID:         2,
		TagType:       1,
		Name:          "Different Tag",
		ImageFilename: sql.NullString{},
	}
	m2_1 := &Tag{
		TagID:         200000,
		TagType:       1,
		Name:          "Different Tag",
		ImageFilename: sql.NullString{},
	}
	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m2_1))

}
