package merger

import (
	"database/sql"
	"sort"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeUserMarkAndBlockRange_without_conflict(t *testing.T) {
	// Merge without conflict
	leftUM := []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			LocationID:   1,
			UserMarkGUID: "DUPLICATE",
		},
		{
			UserMarkID: 2,
			LocationID: 2,
		},
		nil,
		{
			UserMarkID: 4,
			LocationID: 4,
		},
		{},
	}
	leftBR := []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 2,
			UserMarkID:   1,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		nil,
		{
			BlockRangeID: 4,
			UserMarkID:   2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{1, true},
		},
		{
			BlockRangeID: 5,
			UserMarkID:   4,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{2, true},
		},
		{},
	}
	rightUM := []*model.UserMark{
		nil,
		{
			UserMarkID: 1,
			LocationID: 10,
		},
		nil,
		{
			UserMarkID: 3,
			LocationID: 30,
		},
		{
			UserMarkID:   4,
			LocationID:   1,
			UserMarkGUID: "DUPLICATE",
		},
		// Duplicate within right side to check if it is ignored
		{
			UserMarkID: 5,
			LocationID: 30,
		},
	}
	rightBR := []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{30, true},
		},
		{
			BlockRangeID: 1,
			UserMarkID:   4,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 2,
			UserMarkID:   4,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		nil,
		nil,
		{
			BlockRangeID: 4,
			UserMarkID:   3,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 4,
			UserMarkID:   5,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
	}

	expectedUM := []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			LocationID:   1,
			UserMarkGUID: "DUPLICATE",
		},
		{
			UserMarkID: 2,
			LocationID: 2,
		},
		{
			UserMarkID: 3,
			LocationID: 4,
		},
		{
			UserMarkID: 4,
			LocationID: 10,
		},
		{
			UserMarkID: 5,
			LocationID: 30,
		},
		{
			UserMarkID: 6,
			LocationID: 30,
		},
	}
	expectedBR := []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 2,
			UserMarkID:   1,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 3,
			UserMarkID:   2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{1, true},
		},
		{
			BlockRangeID: 4,
			UserMarkID:   3,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{2, true},
		},
		{
			BlockRangeID: 5,
			UserMarkID:   4,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{30, true},
		},
		{
			BlockRangeID: 6,
			UserMarkID:   5,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 7,
			UserMarkID:   6,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
	}
	expectedChanges := IDChanges{
		Left: map[int]int{
			4: 3,
		},
		Right: map[int]int{
			1: 4,
			3: 5,
			4: 1,
			5: 6,
		},
	}

	um, br, changes, err := MergeUserMarkAndBlockRange(leftUM, leftBR, rightUM, rightBR, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedUM, um)
	assert.Equal(t, expectedBR, br)
	assert.Equal(t, expectedChanges, changes)
}

func TestMergeUserMarkAndBlockRange_with_conflict(t *testing.T) {
	// Try merge and find conflict
	leftUM := []*model.UserMark{
		nil,
		{
			UserMarkID: 1,
			LocationID: 1,
		},
		{
			UserMarkID:   2,
			LocationID:   2,
			UserMarkGUID: "DUPLICATE",
		},
		nil,
		{
			UserMarkID: 4,
			LocationID: 4,
		},
		{
			UserMarkID: 5,
			LocationID: 10,
		},
	}
	leftBR := []*model.BlockRange{
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 2,
			UserMarkID:   1,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 3,
			UserMarkID:   2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{1, true},
		},
		{
			BlockRangeID: 4,
			UserMarkID:   4,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{2, true},
		},
		{
			BlockRangeID: 5,
			UserMarkID:   5,
			Identifier:   1,
			StartToken:   sql.NullInt32{29, true},
			EndToken:     sql.NullInt32{35, true},
		},
	}

	rightUM := []*model.UserMark{
		nil,
		{
			UserMarkID: 1,
			LocationID: 10,
		},
		nil,
		{
			UserMarkID: 3,
			LocationID: 30,
		},
		{
			UserMarkID:   4,
			LocationID:   2,
			UserMarkGUID: "DUPLICATEEEEE",
		},
		{
			UserMarkID: 5,
			LocationID: 1,
		},
	}
	rightBR := []*model.BlockRange{
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{30, true},
		},
		{
			BlockRangeID: 2,
			UserMarkID:   3,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 3,
			UserMarkID:   5,
			Identifier:   2,
			StartToken:   sql.NullInt32{3, true},
			EndToken:     sql.NullInt32{7, true},
		},
		{
			BlockRangeID: 4,
			UserMarkID:   4,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{1, true},
		},
	}

	expectedConflicts := []MergeConflict{
		{
			Left: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
					{
						BlockRangeID: 2,
						UserMarkID:   1,
						Identifier:   2,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
				},
			},
			Right: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 5,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 3,
						UserMarkID:   5,
						Identifier:   2,
						StartToken:   sql.NullInt32{3, true},
						EndToken:     sql.NullInt32{7, true},
					},
				},
			},
		},
		{
			Left: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 5,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 5,
						UserMarkID:   5,
						Identifier:   1,
						StartToken:   sql.NullInt32{29, true},
						EndToken:     sql.NullInt32{35, true},
					},
				},
			},
			Right: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{30, true},
					},
				},
			},
		},
	}

	_, _, _, err := MergeUserMarkAndBlockRange(leftUM, leftBR, rightUM, rightBR, nil)
	conflictResult := mergeConflictMapToSliceHelper(err.(MergeConflictError).Conflicts)
	assert.Error(t, err)
	assert.Equal(t, expectedConflicts, conflictResult)

	// Solve conflict
	conflictSolution := map[string]MergeSolution{
		// Merge both markings to one
		"0": {
			Side: LeftSide,
			Solution: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
					{
						BlockRangeID: 2,
						UserMarkID:   1,
						Identifier:   2,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{7, true},
					},
				},
			},
			Discarded: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 5,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 3,
						UserMarkID:   5,
						Identifier:   2,
						StartToken:   sql.NullInt32{3, true},
						EndToken:     sql.NullInt32{7, true},
					},
				},
			},
		},
		"1": {
			Side: RightSide,
			Solution: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{30, true},
					},
				},
			},
			Discarded: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 5,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 5,
						UserMarkID:   5,
						Identifier:   1,
						StartToken:   sql.NullInt32{29, true},
						EndToken:     sql.NullInt32{35, true},
					},
				},
			},
		},
	}

	expectedUM := []*model.UserMark{
		nil,
		{
			UserMarkID: 1,
			LocationID: 1,
		},
		{
			UserMarkID:   2,
			LocationID:   2,
			UserMarkGUID: "DUPLICATE",
		},
		{
			UserMarkID: 3,
			LocationID: 4,
		},
		{
			UserMarkID: 4,
			LocationID: 10,
		},
		{
			UserMarkID: 5,
			LocationID: 30,
		},
	}
	expectedBR := []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
		{
			BlockRangeID: 2,
			UserMarkID:   1,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{7, true},
		},
		{
			BlockRangeID: 3,
			UserMarkID:   2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{1, true},
		},
		{
			BlockRangeID: 4,
			UserMarkID:   3,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{2, true},
		},
		{
			BlockRangeID: 5,
			UserMarkID:   4,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{30, true},
		},
		{
			BlockRangeID: 6,
			UserMarkID:   5,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{5, true},
		},
	}
	expectedChanges := IDChanges{
		Left: map[int]int{
			4: 3,
			5: 4,
		},
		Right: map[int]int{
			1: 4,
			3: 5,
			4: 2,
			5: 1,
		},
	}

	um, br, changes, err := MergeUserMarkAndBlockRange(leftUM, leftBR, rightUM, rightBR, conflictSolution)
	assert.NoError(t, err)
	assert.Equal(t, expectedUM, um)
	assert.Equal(t, expectedBR, br)
	assert.Equal(t, expectedChanges, changes)
}

func Test_mergeUserMarkBlockRange_without_conflict(t *testing.T) {
	// Merge without conflict
	left := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
				{
					BlockRangeID: 2,
					UserMarkID:   1,
					Identifier:   2,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   2,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 4,
				LocationID: 4,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   4,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{2, true},
				},
			},
		},
	}

	right := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 10,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{30, true},
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 3,
				LocationID: 30,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 2,
					UserMarkID:   3,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
	}

	expectedResult := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
				{
					BlockRangeID: 2,
					UserMarkID:   1,
					Identifier:   2,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   2,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 3,
				LocationID: 4,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   4,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{2, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 4,
				LocationID: 10,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{30, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 5,
				LocationID: 30,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 2,
					UserMarkID:   3,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
	}
	expectedChanges := IDChanges{
		Left: map[int]int{
			4: 3,
		},
		Right: map[int]int{
			1: 4,
			3: 5,
		},
	}

	result, changes, err := mergeUMBR(left, right, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedChanges, changes)
}

func Test_mergeUserMarkBlockRange_with_conflict(t *testing.T) {
	// Try merge and find conflict
	left := []*model.UserMarkBlockRange{
		nil,
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   2,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
				{
					BlockRangeID: 2,
					UserMarkID:   2,
					Identifier:   2,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 3,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   3,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 5,
				LocationID: 4,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   5,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{2, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 6,
				LocationID: 10,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   6,
					Identifier:   1,
					StartToken:   sql.NullInt32{29, true},
					EndToken:     sql.NullInt32{35, true},
				},
			},
		},
	}

	right := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 10,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{30, true},
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 3,
				LocationID: 30,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 2,
					UserMarkID:   3,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 4,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   4,
					Identifier:   2,
					StartToken:   sql.NullInt32{3, true},
					EndToken:     sql.NullInt32{7, true},
				},
			},
		},
	}

	expectedConflicts := []MergeConflict{
		{
			Left: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 2,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   2,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
					{
						BlockRangeID: 2,
						UserMarkID:   2,
						Identifier:   2,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
				},
			},
			Right: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 4,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   4,
						Identifier:   2,
						StartToken:   sql.NullInt32{3, true},
						EndToken:     sql.NullInt32{7, true},
					},
				},
			},
		},
		{
			Left: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 6,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   6,
						Identifier:   1,
						StartToken:   sql.NullInt32{29, true},
						EndToken:     sql.NullInt32{35, true},
					},
				},
			},
			Right: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{30, true},
					},
				},
			},
		},
	}

	result, _, err := mergeUMBR(left, right, nil)
	conflictResult := mergeConflictMapToSliceHelper(err.(MergeConflictError).Conflicts)
	assert.Empty(t, result)
	assert.Error(t, err)
	assert.Equal(t, expectedConflicts, conflictResult)

	// Solve conflict
	conflictSolution := map[string]MergeSolution{
		// Merge both markings to one
		"0": {
			Side: LeftSide,
			Solution: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 2,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   2,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
					{
						BlockRangeID: 2,
						UserMarkID:   2,
						Identifier:   2,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{7, true},
					},
				},
			},
			Discarded: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 4,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   4,
						Identifier:   2,
						StartToken:   sql.NullInt32{3, true},
						EndToken:     sql.NullInt32{7, true},
					},
				},
			},
		},
		"1": {
			Side: RightSide,
			Solution: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{30, true},
					},
				},
			},
			Discarded: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 6,
					LocationID: 10,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   6,
						Identifier:   1,
						StartToken:   sql.NullInt32{29, true},
						EndToken:     sql.NullInt32{35, true},
					},
				},
			},
		},
	}

	expectedResult := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   2, // splitUserMarkBlockRange will fix this error in a later step
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
				{
					BlockRangeID: 2,
					UserMarkID:   2,
					Identifier:   2,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{7, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   3,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 3,
				LocationID: 4,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   5,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{2, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 4,
				LocationID: 10,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{30, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 5,
				LocationID: 30,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 2,
					UserMarkID:   3,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
	}
	expectedChanges := IDChanges{
		Left: map[int]int{
			2: 1,
			3: 2,
			5: 3,
			6: 4,
		},
		Right: map[int]int{
			1: 4,
			3: 5,
			4: 1,
		},
	}

	result, changes, err := mergeUMBR(left, right, conflictSolution)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedChanges, changes)
}

// mergeConflictMapToSliceHelper is a helper function that converts a mergeConflict map
// to a sorted slice. This makes testing reliable, as we are able to trust
// the order of a map.
func mergeConflictMapToSliceHelper(mp map[string]MergeConflict) []MergeConflict {
	result := []MergeConflict{}
	for _, entry := range mp {
		result = append(result, entry)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Left.ID() < result[j].Left.ID()
	})

	return result
}

func Test_replaceUMBRConflictsWithSolution(t *testing.T) {
	left := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 3,
				LocationID: 3,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{2, true},
				},
			},
		},
	}

	right := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 3,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
	}

	conflictSolution := map[string]MergeSolution{
		"0": {
			Side: LeftSide,
			Solution: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{1, true},
					},
				},
			},
			Discarded: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 1,
					LocationID: 1,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{1, true},
					},
				},
			},
		},
		"1": {
			Side: RightSide,
			Solution: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 2,
					LocationID: 3,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{5, true},
					},
				},
			},
			Discarded: &model.UserMarkBlockRange{
				UserMark: &model.UserMark{
					UserMarkID: 3,
					LocationID: 3,
				},
				BlockRanges: []*model.BlockRange{
					{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
						StartToken:   sql.NullInt32{0, true},
						EndToken:     sql.NullInt32{2, true},
					},
				},
			},
		},
	}

	expectedLeft := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{1, true},
				},
			},
		},
		nil,
	}
	expectedRight := []*model.UserMarkBlockRange{
		nil,
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 3,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
					StartToken:   sql.NullInt32{0, true},
					EndToken:     sql.NullInt32{5, true},
				},
			},
		},
	}
	expectedChanges := IDChanges{
		Left: map[int]int{
			3: 2,
		},
		Right: map[int]int{
			1: 1,
		},
	}
	expectedInvertedChanges := IDChanges{
		Left: map[int]int{
			2: 3,
		},
		Right: map[int]int{
			1: 1,
		},
	}

	changes, invertedChanges := replaceUMBRConflictsWithSolution(&left, &right, conflictSolution)
	assert.Equal(t, expectedLeft, left)
	assert.Equal(t, expectedRight, right)
	assert.Equal(t, expectedChanges, changes)
	assert.Equal(t, expectedInvertedChanges, invertedChanges)
}

func Test_ingestUMBR(t *testing.T) {
	left := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
				},
				{
					BlockRangeID: 2,
					UserMarkID:   1,
					Identifier:   1,
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 3,
					UserMarkID:   2,
					Identifier:   1,
				},
				{
					BlockRangeID: 4,
					UserMarkID:   2,
					Identifier:   2,
				},
			},
		},
		nil,
	}

	right := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
					Identifier:   1,
				},
				{
					BlockRangeID: 2,
					UserMarkID:   1,
					Identifier:   1,
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
				LocationID: 20,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 3,
					UserMarkID:   2,
					Identifier:   1,
				},
				{
					BlockRangeID: 4,
					UserMarkID:   2,
					Identifier:   2,
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 4,
				LocationID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 4,
					UserMarkID:   4,
					Identifier:   1,
				},
				{
					BlockRangeID: 5,
					UserMarkID:   4,
					Identifier:   10,
				},
			},
		},
	}

	// Map[LocationID]map[Identifier][]*model.BlockRange
	expectedResult := map[int]map[int][]brFrom{
		1: {
			1: []brFrom{
				{
					side: LeftSide,
					br: &model.BlockRange{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
					},
				},
				{
					side: LeftSide,
					br: &model.BlockRange{
						BlockRangeID: 2,
						UserMarkID:   1,
						Identifier:   1,
					},
				},
				{
					side: RightSide,
					br: &model.BlockRange{
						BlockRangeID: 1,
						UserMarkID:   1,
						Identifier:   1,
					},
				},
				{
					side: RightSide,
					br: &model.BlockRange{
						BlockRangeID: 2,
						UserMarkID:   1,
						Identifier:   1,
					},
				},
				{
					side: RightSide,
					br: &model.BlockRange{
						BlockRangeID: 4,
						UserMarkID:   4,
						Identifier:   1,
					},
				},
			},
			10: []brFrom{
				{
					side: RightSide,
					br: &model.BlockRange{
						BlockRangeID: 5,
						UserMarkID:   4,
						Identifier:   10,
					},
				},
			},
		},
		2: {
			1: []brFrom{
				{
					side: LeftSide,
					br: &model.BlockRange{
						BlockRangeID: 3,
						UserMarkID:   2,
						Identifier:   1,
					},
				},
			},
			2: []brFrom{
				{
					side: LeftSide,
					br: &model.BlockRange{
						BlockRangeID: 4,
						UserMarkID:   2,
						Identifier:   2,
					},
				},
			},
		},
		20: {
			1: []brFrom{
				{
					side: RightSide,
					br: &model.BlockRange{
						BlockRangeID: 3,
						UserMarkID:   2,
						Identifier:   1,
					},
				},
			},
			2: []brFrom{
				{
					side: RightSide,
					br: &model.BlockRange{
						BlockRangeID: 4,
						UserMarkID:   2,
						Identifier:   2,
					},
				},
			},
		},
	}

	result := ingestUMBR(left, right)
	assert.Equal(t, expectedResult, result)
}

func Test_joinToUserMarkBlockRange(t *testing.T) {
	userMarks := []*model.UserMark{
		nil,
		{
			UserMarkID: 1,
		},
		{
			UserMarkID: 2,
		},
		nil,
		{
			UserMarkID: 4,
		},
		{},
	}
	blockRanges := []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			UserMarkID:   1,
		},
		{
			BlockRangeID: 2,
			UserMarkID:   1,
		},
		{
			BlockRangeID: 3,
			UserMarkID:   2,
		},
		nil,
		{
			BlockRangeID: 4,
			UserMarkID:   2,
		},
		{},
		{
			BlockRangeID: 5,
			UserMarkID:   12345,
		},
	}

	expectedResult := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 1,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 1,
					UserMarkID:   1,
				},
				{
					BlockRangeID: 2,
					UserMarkID:   1,
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID: 2,
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 3,
					UserMarkID:   2,
				},
				{
					BlockRangeID: 4,
					UserMarkID:   2,
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID: 4,
			},
			BlockRanges: []*model.BlockRange{},
		},
		nil,
	}

	result := joinToUserMarkBlockRange(userMarks, blockRanges)
	assert.Equal(t, expectedResult, result)

}

func Test_splitUserMarkBlockRange(t *testing.T) {
	umbr := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID:   1,
				UserMarkGUID: "#1",
			},
			BlockRanges: []*model.BlockRange{
				nil,
				{
					BlockRangeID: 4,
					UserMarkID:   20,
					Identifier:   1,
				},
				{
					BlockRangeID: 2,
					UserMarkID:   1,
					Identifier:   2,
				},
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID:   2,
				UserMarkGUID: "#2",
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 3,
					UserMarkID:   5,
					Identifier:   3,
				},
				{},
				{
					BlockRangeID: 3,
					UserMarkID:   5,
					Identifier:   4,
				},
			},
		},
		nil,
		{
			UserMark: &model.UserMark{
				UserMarkID:   4,
				UserMarkGUID: "#4",
			},
			BlockRanges: []*model.BlockRange{},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID:   5,
				UserMarkGUID: "#5",
			},
		},
		{
			UserMark: &model.UserMark{
				UserMarkID:   6,
				UserMarkGUID: "#6",
			},
			BlockRanges: []*model.BlockRange{
				{
					BlockRangeID: 3,
					UserMarkID:   6,
					Identifier:   5,
				},
				nil,
				{
					BlockRangeID: 3,
					UserMarkID:   6,
					Identifier:   6,
				},
			},
		},
	}

	expectedUM := []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			UserMarkGUID: "#1",
		},
		{
			UserMarkID:   2,
			UserMarkGUID: "#2",
		},
		nil,
		{
			UserMarkID:   4,
			UserMarkGUID: "#4",
		},
		{
			UserMarkID:   5,
			UserMarkGUID: "#5",
		},
		{
			UserMarkID:   6,
			UserMarkGUID: "#6",
		},
	}

	expectedBR := []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			UserMarkID:   1,
			Identifier:   1,
		},
		{
			BlockRangeID: 2,
			UserMarkID:   1,
			Identifier:   2,
		},
		{
			BlockRangeID: 3,
			UserMarkID:   2,
			Identifier:   3,
		},
		{
			BlockRangeID: 4,
			UserMarkID:   2,
			Identifier:   4,
		},
		{
			BlockRangeID: 5,
			UserMarkID:   6,
			Identifier:   5,
		},
		{
			BlockRangeID: 6,
			UserMarkID:   6,
			Identifier:   6,
		},
	}

	um, br := splitUserMarkBlockRange(umbr)
	assert.Equal(t, expectedUM, um)
	assert.Equal(t, expectedBR, br)

	umbr = []*model.UserMarkBlockRange{}
	um, br = splitUserMarkBlockRange(umbr)
	assert.Empty(t, um)
	assert.Empty(t, br)
}

func Test_estimateLocationCount(t *testing.T) {
	left := []*model.UserMarkBlockRange{
		nil,
		{
			UserMark: &model.UserMark{
				LocationID: 100,
			},
		},
		{
			UserMark: &model.UserMark{
				LocationID: 1000,
			},
		},
		{
			UserMark: &model.UserMark{
				LocationID: 1234,
			},
		},
	}

	right := []*model.UserMarkBlockRange{
		nil,
		{UserMark: &model.UserMark{
			LocationID: 10,
		},
		},
		{
			UserMark: &model.UserMark{
				LocationID: 1500,
			},
		},
		{
			UserMark: &model.UserMark{
				LocationID: 2234,
			},
		},
	}

	assert.Equal(t, 3468, estimateLocationCount(left, right))
	assert.Equal(t, 100, estimateLocationCount([]*model.UserMarkBlockRange{}, []*model.UserMarkBlockRange{{UserMark: &model.UserMark{LocationID: 100}}}))
	assert.Equal(t, 101, estimateLocationCount([]*model.UserMarkBlockRange{{UserMark: &model.UserMark{LocationID: 101}}}, []*model.UserMarkBlockRange{}))
	assert.Equal(t, 0, estimateLocationCount([]*model.UserMarkBlockRange{}, []*model.UserMarkBlockRange{}))
	assert.Equal(t, 0, estimateLocationCount([]*model.UserMarkBlockRange{nil}, []*model.UserMarkBlockRange{nil}))
}
