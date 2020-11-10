package gomobile

import (
	"database/sql"
	"encoding/json"
	"sort"
	"testing"

	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeConflictsWrapper_InitDBWrapper(t *testing.T) {
	dbw := &DatabaseWrapper{
		left: &model.Database{},
	}

	mcw := MergeConflictsWrapper{}
	mcw.InitDBWrapper(dbw)

	assert.Same(t, dbw, mcw.DBWrapper)
}

func TestMergeConflictsWrapper_addConflicts(t *testing.T) {
	conflicts := map[string]merger.MergeConflict{
		"1": {
			Left: &model.Bookmark{
				Title: "1Left",
			},
			Right: &model.Bookmark{
				Title: "1Right",
			},
		},
		"2": {
			Left: &model.Tag{
				Name: "2Left",
			},
			Right: &model.Tag{
				Name: "2Right",
			},
		},
	}

	mcw := MergeConflictsWrapper{}
	mcw.addConflicts(conflicts)
	sort.Strings(mcw.conflictKeys)
	assert.Equal(t, conflicts, mcw.conflicts)
	assert.Equal(t, []string{"1", "2"}, mcw.conflictKeys)

	conflicts["3"] = merger.MergeConflict{
		Left: &model.Tag{
			Name: "2Left",
		},
		Right: &model.Tag{
			Name: "2Right",
		},
	}

	mcw.addConflicts(conflicts)
	assert.Equal(t, conflicts, mcw.conflicts)
	assert.Equal(t, []string{"1", "2", "3"}, mcw.conflictKeys)
}

func TestMergeConflictsWrapper_ConflictsLen(t *testing.T) {
	mcw := MergeConflictsWrapper{}
	assert.Equal(t, 0, mcw.ConflictsLen())

	mcw.conflicts = map[string]merger.MergeConflict{"1": {}, "2": {}, "3": {}}
	assert.Equal(t, 3, mcw.ConflictsLen())
}

func TestMergeConflictsWrapper_SolutionsLen(t *testing.T) {
	mcw := MergeConflictsWrapper{}
	assert.Equal(t, 0, mcw.SolutionsLen())

	mcw.solutions = map[string]merger.MergeSolution{"1": {}, "2": {}, "3": {}}
	assert.Equal(t, 3, mcw.SolutionsLen())
}

func TestMergeConflictsWrapper_GetNextConflictIndex(t *testing.T) {
	conflicts := map[string]merger.MergeConflict{
		"1": {
			Left: &model.Bookmark{
				Title: "1Left",
			},
			Right: &model.Bookmark{
				Title: "1Right",
			},
		},
		"2": {
			Left: &model.Tag{
				Name: "2Left",
			},
			Right: &model.Tag{
				Name: "2Right",
			},
		},
	}

	mcw := MergeConflictsWrapper{}
	mcw.addConflicts(conflicts)

	assert.Equal(t, 0, mcw.GetNextConflictIndex())
	assert.Error(t, mcw.SolveConflict(1, "leftSide"))
	assert.NoError(t, mcw.SolveConflict(0, "leftSide"))
	assert.NoError(t, mcw.SolveConflict(0, "leftSide")) // Solving the same conflict should not increase counter
	assert.Equal(t, 1, mcw.GetNextConflictIndex())
	assert.NoError(t, mcw.SolveConflict(1, "leftSide"))
	assert.Equal(t, -1, mcw.GetNextConflictIndex())
}

func TestMergeConflictsWrapper_GetConflict(t *testing.T) {
	db := &model.Database{
		Location: []*model.Location{
			nil,
			{
				LocationID: 1,
				Title:      sql.NullString{"Location-Title", true},
			},
		},
	}

	mcw := MergeConflictsWrapper{
		DBWrapper: &DatabaseWrapper{
			merged: db,
		},
		conflicts: map[string]merger.MergeConflict{
			"1": {
				Left: &model.Bookmark{
					LocationID: 1,
					Title:      "1Left",
				},
				Right: &model.Bookmark{
					LocationID: 1,
					Title:      "1Right",
				},
			},
			"2": {
				Left: &model.Tag{
					Name: "2Left",
				},
				Right: &model.Tag{
					Name: "2Right",
				},
			},
		},
		conflictKeys: []string{"1", "2"},
	}

	conflict, err := mcw.GetConflict(0)
	assert.NoError(t, err)
	assert.Equal(t,
		jsonMarhshalIgnoreErr(modelRelatedTuple{
			Model:   mcw.conflicts["1"].Left,
			Related: model.Related{Location: db.Location[1], PublicationLocation: db.Location[1]},
		}),
		conflict.Left)
	assert.Equal(t,
		jsonMarhshalIgnoreErr(modelRelatedTuple{
			Model:   mcw.conflicts["1"].Right,
			Related: model.Related{Location: db.Location[1], PublicationLocation: db.Location[1]},
		}),
		conflict.Right)

	conflict, err = mcw.GetConflict(1)
	assert.NoError(t, err)
	assert.Equal(t,
		jsonMarhshalIgnoreErr(modelRelatedTuple{
			Model:   mcw.conflicts["2"].Left,
			Related: model.Related{},
		}),
		conflict.Left)
	assert.Equal(t,
		jsonMarhshalIgnoreErr(modelRelatedTuple{
			Model:   mcw.conflicts["2"].Right,
			Related: model.Related{},
		}),
		conflict.Right)

	mcw.DBWrapper = nil
	assert.Equal(t,
		jsonMarhshalIgnoreErr(modelRelatedTuple{
			Model:   mcw.conflicts["2"].Right,
			Related: model.Related{},
		}),
		conflict.Right)

	_, err = mcw.GetConflict(3)
	assert.Error(t, err)

	mcw.conflicts = nil
	_, err = mcw.GetConflict(1)
	assert.Error(t, err)
}

func jsonMarhshalIgnoreErr(m interface{}) string {
	result, _ := json.Marshal(m)
	return string(result)
}

func TestMergeConflictsWrapper_SolveConflict(t *testing.T) {
	mcw := MergeConflictsWrapper{}
	assert.EqualError(t, mcw.SolveConflict(0, "leftSide"), "There are no conflicts")

	mcw = MergeConflictsWrapper{
		conflicts: map[string]merger.MergeConflict{
			"1": {
				Left: &model.Bookmark{
					Title: "1Left",
				},
				Right: &model.Bookmark{
					Title: "1Right",
				},
			},
			"2": {
				Left: &model.Tag{
					Name: "2Left",
				},
				Right: &model.Tag{
					Name: "2Right",
				},
			},
			"3": {
				Left: &model.Tag{
					Name: "3Left",
				},
				Right: &model.Tag{
					Name: "3Right",
				},
			},
		},
		conflictKeys: []string{"1", "2", "3"},
	}
	assert.EqualError(t, mcw.SolveConflict(3, "leftSide"), "Conflict with index 3 does not exist. Length=3")
	assert.EqualError(t, mcw.SolveConflict(1, "leftSide"), "Index is higher than NextConflictIndex: 1 > 0. The conflicts before have to be solved first")
	assert.EqualError(t, mcw.SolveConflict(0, "wrongSide"), "Side wrongSide is not valid")

	assert.NoError(t, mcw.SolveConflict(0, "rightSide"))
	assert.Equal(t,
		merger.MergeSolution{
			Side:      merger.RightSide,
			Solution:  mcw.conflicts["1"].Right,
			Discarded: mcw.conflicts["1"].Left,
		},
		mcw.solutions["1"])
	assert.Equal(t, 1, mcw.GetNextConflictIndex())
	assert.NoError(t, mcw.SolveConflict(0, "rightSide"))
	assert.Equal(t, 1, mcw.GetNextConflictIndex())

	assert.NoError(t, mcw.SolveConflict(1, "leftSide"))
	assert.Equal(t,
		merger.MergeSolution{
			Side:      merger.LeftSide,
			Solution:  mcw.conflicts["2"].Left,
			Discarded: mcw.conflicts["2"].Right,
		},
		mcw.solutions["2"])
	assert.Equal(t, 2, mcw.GetNextConflictIndex())

	assert.NoError(t, mcw.SolveConflict(2, "leftSide"))
	assert.Equal(t, -1, mcw.GetNextConflictIndex())

	// Change previous solution to different one
	assert.NoError(t, mcw.SolveConflict(0, "leftSide"))
	assert.Equal(t,
		merger.MergeSolution{
			Side:      merger.LeftSide,
			Solution:  mcw.conflicts["1"].Left,
			Discarded: mcw.conflicts["1"].Right,
		},
		mcw.solutions["1"])
	assert.Equal(t, -1, mcw.GetNextConflictIndex())
}
