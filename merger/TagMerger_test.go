package merger

import (
	"database/sql"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeTags(t *testing.T) {
	// Successful merge
	left := []*model.Tag{
		{
			TagID:         1,
			TagType:       1,
			Name:          "A tag on the left",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         2,
			TagType:       1,
			Name:          "Another tag on the left",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         3,
			TagType:       1,
			Name:          "A duplicate tag",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         4,
			TagType:       1,
			Name:          "A duplicate tag with imageFilename",
			ImageFilename: sql.NullString{String: "aFileName.jpg", Valid: true},
		},
		nil,
	}
	right := []*model.Tag{
		{
			TagID:         1,
			TagType:       1,
			Name:          "A tag on the right",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         2,
			TagType:       1,
			Name:          "A duplicate tag",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         3,
			TagType:       1,
			Name:          "A duplicate tag with imageFilename",
			ImageFilename: sql.NullString{String: "aFileName.jpg", Valid: true},
		},
		{
			TagID:         4,
			TagType:       1,
			Name:          "One more tag only on the right",
			ImageFilename: sql.NullString{},
		},
	}

	expectedResult := []*model.Tag{
		nil,
		{
			TagID:         1,
			TagType:       1,
			Name:          "A tag on the left",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         2,
			TagType:       1,
			Name:          "A tag on the right",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         3,
			TagType:       1,
			Name:          "Another tag on the left",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         4,
			TagType:       1,
			Name:          "A duplicate tag",
			ImageFilename: sql.NullString{},
		},
		{
			TagID:         5,
			TagType:       1,
			Name:          "A duplicate tag with imageFilename",
			ImageFilename: sql.NullString{String: "aFileName.jpg", Valid: true},
		},
		{
			TagID:         6,
			TagType:       1,
			Name:          "One more tag only on the right",
			ImageFilename: sql.NullString{},
		},
	}

	expectedChanges := IDChanges{
		Left: map[int]int{
			2: 3,
			3: 4,
			4: 5,
		},
		Right: map[int]int{
			1: 2,
			2: 4,
			3: 5,
			4: 6,
		},
	}

	result, changes, err := MergeTags(left, right, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedChanges, changes)
}
