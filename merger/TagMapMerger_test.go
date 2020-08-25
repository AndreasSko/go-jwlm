package merger

import (
	"database/sql"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeTagMaps(t *testing.T) {
	left := []*model.TagMap{
		{
			TagMapID:       1,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 1, Valid: true},
			TagID:          1,
			Position:       0,
		},
		{
			TagMapID:       2,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          2,
			Position:       0,
		},
		{
			TagMapID:       3,
			PlaylistItemID: sql.NullInt32{Int32: 1, Valid: true},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{},
			TagID:          5,
			Position:       0,
		},
		{
			TagMapID:       4,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          3,
			Position:       1,
		},
		{
			TagMapID:       5,
			PlaylistItemID: sql.NullInt32{Int32: 1, Valid: true},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{},
			TagID:          4,
			Position:       2,
		},
		nil,
	}

	right := []*model.TagMap{
		{
			TagMapID:       1,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          4,
			Position:       0,
		},
		{
			TagMapID:       2,
			PlaylistItemID: sql.NullInt32{Int32: 1, Valid: true},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{},
			TagID:          6,
			Position:       0,
		},
		// Duplicate of left TagMapID 1
		{
			TagMapID:       3,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 1, Valid: true},
			TagID:          1,
			Position:       0,
		},
		{
			TagMapID:       4,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          3,
			Position:       2,
		},
	}

	expectedResult := []*model.TagMap{
		nil,
		{
			TagMapID:       1,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 1, Valid: true},
			TagID:          1,
			Position:       0,
		},
		{
			TagMapID:       2,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          2,
			Position:       0,
		},
		{
			TagMapID:       3,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          4,
			Position:       1,
		},
		{
			TagMapID:       4,
			PlaylistItemID: sql.NullInt32{},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{Int32: 2, Valid: true},
			TagID:          3,
			Position:       2,
		},
		{
			TagMapID:       5,
			PlaylistItemID: sql.NullInt32{Int32: 1, Valid: true},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{},
			TagID:          5,
			Position:       0,
		},
		{
			TagMapID:       6,
			PlaylistItemID: sql.NullInt32{Int32: 1, Valid: true},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{},
			TagID:          6,
			Position:       1,
		},
		{
			TagMapID:       7,
			PlaylistItemID: sql.NullInt32{Int32: 1, Valid: true},
			LocationID:     sql.NullInt32{},
			NoteID:         sql.NullInt32{},
			TagID:          4,
			Position:       2,
		},
	}

	result, _, err := MergeTagMaps(left, right, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}
