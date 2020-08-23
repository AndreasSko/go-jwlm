package model

import (
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
