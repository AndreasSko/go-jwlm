package merger

import (
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestSolveConflictByChoosingX(t *testing.T) {
	conflicts := map[string]MergeConflict{
		"bookmarkCollision": {
			Left: &model.Bookmark{
				Title: "LeftBookmark",
			},
			Right: &model.Bookmark{
				Title: "RightBookmark",
			},
		},
		"noteCollision": {
			Left: &model.Note{
				GUID: "LeftNote",
			},
			Right: &model.Note{
				GUID: "RightNote",
			},
		},
	}

	// Choose left
	expectedResult := map[string]MergeSolution{
		"bookmarkCollision": {
			Side: LeftSide,
			Solution: &model.Bookmark{
				Title: "LeftBookmark",
			},
			Discarded: &model.Bookmark{
				Title: "RightBookmark",
			},
		},
		"noteCollision": {
			Side: LeftSide,
			Solution: &model.Note{
				GUID: "LeftNote",
			},
			Discarded: &model.Note{
				GUID: "RightNote",
			},
		},
	}

	result, err := SolveConflictByChoosingLeft(conflicts)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	// Choose right
	expectedResult = map[string]MergeSolution{
		"bookmarkCollision": {
			Side: RightSide,
			Solution: &model.Bookmark{
				Title: "RightBookmark",
			},
			Discarded: &model.Bookmark{
				Title: "LeftBookmark",
			},
		},
		"noteCollision": {
			Side: RightSide,
			Solution: &model.Note{
				GUID: "RightNote",
			},
			Discarded: &model.Note{
				GUID: "LeftNote",
			},
		},
	}

	result, err = SolveConflictByChoosingRight(conflicts)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSolveConflictByChoosingNewest(t *testing.T) {
	conflicts := map[string]MergeConflict{
		"leftNewer": {
			Left: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
			},
			Right: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T12:00:00+00:00",
			},
		},
		"rightNewer": {
			Left: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
			},
			Right: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T13:01:00+00:00",
			},
		},
	}

	expectedResult := map[string]MergeSolution{
		"leftNewer": {
			Side: LeftSide,
			Solution: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
			},
			Discarded: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T12:00:00+00:00",
			},
		},
		"rightNewer": {
			Side: RightSide,
			Solution: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T13:01:00+00:00",
			},
			Discarded: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
			},
		},
	}

	result, err := SolveConflictByChoosingNewest(conflicts)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	conflicts = map[string]MergeConflict{
		"bookmarkCollision": {
			Left: &model.Bookmark{
				Title: "LeftBookmark",
			},
			Right: &model.Bookmark{
				Title: "RightBookmark",
			},
		},
	}
	_, err = SolveConflictByChoosingNewest(conflicts)
	assert.Error(t, err)
}
