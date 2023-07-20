package merger

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestAutoResolveConflicts(t *testing.T) {
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

	result, err := AutoResolveConflicts(conflicts, "chooseNewest")
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = AutoResolveConflicts(conflicts, "wrongResolver")
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = AutoResolveConflicts(conflicts, "")
	assert.NoError(t, err)
	assert.Nil(t, result)

	result, err = AutoResolveConflicts(conflicts, "chooseLeft")
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

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
				Created:      "2020-09-15T12:00:00+00:00",
			},
			Right: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T12:00:00+00:00",
				Created:      "2020-09-15T12:00:00+00:00",
			},
		},
		"rightNewer": {
			Left: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
				Created:      "2020-09-15T13:00:00+00:00",
			},
			Right: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T13:01:00+00:00",
				Created:      "2020-09-15T13:00:00+00:00",
			},
		},
	}

	expectedResult := map[string]MergeSolution{
		"leftNewer": {
			Side: LeftSide,
			Solution: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
				Created:      "2020-09-15T12:00:00+00:00",
			},
			Discarded: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T12:00:00+00:00",
				Created:      "2020-09-15T12:00:00+00:00",
			},
		},
		"rightNewer": {
			Side: RightSide,
			Solution: &model.Note{
				GUID:         "Right",
				LastModified: "2020-09-15T13:01:00+00:00",
				Created:      "2020-09-15T13:00:00+00:00",
			},
			Discarded: &model.Note{
				GUID:         "Left",
				LastModified: "2020-09-15T13:00:00+00:00",
				Created:      "2020-09-15T13:00:00+00:00",
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

func Test_parseResolver(t *testing.T) {
	resolver, err := parseResolver("")
	assert.NoError(t, err)
	assert.Nil(t, resolver)

	// https://github.com/stretchr/testify/issues/182#issuecomment-495359313
	resolver, err = parseResolver("chooseLeft")
	assert.NoError(t, err)
	assert.Equal(t,
		"github.com/AndreasSko/go-jwlm/merger.SolveConflictByChoosingLeft",
		runtime.FuncForPC(reflect.ValueOf(resolver).Pointer()).Name())

	resolver, err = parseResolver("chooseRight")
	assert.NoError(t, err)
	assert.Equal(t,
		"github.com/AndreasSko/go-jwlm/merger.SolveConflictByChoosingRight",
		runtime.FuncForPC(reflect.ValueOf(resolver).Pointer()).Name())

	resolver, err = parseResolver("chooseNewest")
	assert.NoError(t, err)
	assert.Equal(t,
		"github.com/AndreasSko/go-jwlm/merger.SolveConflictByChoosingNewest",
		runtime.FuncForPC(reflect.ValueOf(resolver).Pointer()).Name())

	resolver, err = parseResolver("nonexistent")
	assert.EqualError(t, err, "nonexistent is not a valid conflict resolver. Can be 'chooseNewest', 'chooseLeft', or 'chooseRight'")
	assert.Nil(t, resolver)
}
