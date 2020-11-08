package gomobile

import (
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseWrapper_Stats(t *testing.T) {
	db := &model.Database{
		BlockRange: []*model.BlockRange{nil},
		Bookmark:   []*model.Bookmark{nil, nil},
		Location:   []*model.Location{nil, nil, nil},
		Note:       []*model.Note{nil, nil, nil, nil},
		Tag:        []*model.Tag{nil, nil, nil, nil, nil},
		TagMap:     []*model.TagMap{nil, nil, nil, nil, nil, nil},
		UserMark:   []*model.UserMark{nil, nil, nil, nil, nil, nil, nil},
	}
	dbw := &DatabaseWrapper{
		left:   db,
		right:  db,
		merged: db,
	}

	expected := &DatabaseStats{
		BlockRange: 1,
		Bookmark:   2,
		Location:   3,
		Note:       4,
		Tag:        5,
		TagMap:     6,
		UserMark:   7,
	}

	assert.Equal(t, expected, dbw.Stats("leftSide"))
	assert.Equal(t, expected, dbw.Stats("rightSide"))
	assert.Equal(t, expected, dbw.Stats("mergeSide"))
	assert.Equal(t, &DatabaseStats{}, dbw.Stats("wrongSide"))
}
