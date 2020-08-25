package merger

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func Test_MergeLocations(t *testing.T) {
	left := []*model.Location{
		nil,
		{
			LocationID:     1,
			BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "", Valid: true},
		},
		{
			LocationID:     2,
			BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 4, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ThisTitleShouldStay", Valid: true},
		},
		{
			LocationID:     3,
			BookNumber:     sql.NullInt32{Int32: 5, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 8, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "", Valid: true},
		},
		{
			LocationID:     4,
			BookNumber:     sql.NullInt32{Int32: 10, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 8, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomeTitle", Valid: true},
		},
		{
			LocationID:     5,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 2020401, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "w", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "AWTTitle", Valid: true},
		},
		{
			LocationID:     6,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 2020501, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "w", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomethingDuplicate", Valid: true},
		},
	}

	right := []*model.Location{
		nil,
		{
			LocationID:     1,
			BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "", Valid: true},
		},
		{
			LocationID:     2,
			BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 4, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ThisTitleShouldStay", Valid: true},
		},
		{
			LocationID:     3,
			BookNumber:     sql.NullInt32{Int32: 5, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 8, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ThisTitleShouldStay", Valid: true},
		},
		{
			LocationID:     4,
			BookNumber:     sql.NullInt32{Int32: 6, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 5, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomeOtherTitle", Valid: true},
		},
		{
			LocationID:     5,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 123456789, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "CO-pgm20", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ATitle", Valid: true},
		},
		{
			LocationID:     6,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 123456790, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "CO-pgm20", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "AATitle", Valid: true},
		},
		{
			LocationID:     7,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 2020501, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "w", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomethingDuplicate", Valid: true},
		},
	}

	expectedResult := []*model.Location{
		nil,
		{
			LocationID:     1,
			BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "", Valid: true},
		},
		{
			LocationID:     2,
			BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 4, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ThisTitleShouldStay", Valid: true},
		},
		{
			LocationID:     3,
			BookNumber:     sql.NullInt32{Int32: 5, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 8, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ThisTitleShouldStay", Valid: true},
		},
		{
			LocationID:     4,
			BookNumber:     sql.NullInt32{Int32: 10, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 8, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomeTitle", Valid: true},
		},
		{
			LocationID:     5,
			BookNumber:     sql.NullInt32{Int32: 6, Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: 5, Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomeOtherTitle", Valid: true},
		},
		{
			LocationID:     6,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 2020401, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "w", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "AWTTitle", Valid: true},
		},
		{
			LocationID:     7,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 123456789, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "CO-pgm20", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "ATitle", Valid: true},
		},
		{
			LocationID:     8,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 2020501, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "w", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "SomethingDuplicate", Valid: true},
		},
		{
			LocationID:     9,
			BookNumber:     sql.NullInt32{},
			ChapterNumber:  sql.NullInt32{},
			DocumentID:     sql.NullInt32{Int32: 123456790, Valid: true},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "CO-pgm20", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: "AATitle", Valid: true},
		},
	}

	expectedChanges := IDChanges{
		Left: map[int]int{
			5: 6,
			6: 8,
		},
		Right: map[int]int{
			4: 5,
			5: 7,
			6: 9,
			7: 8,
		},
	}

	result, changes, err := MergeLocations(left, right)

	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedChanges, changes)
	assert.NoError(t, err)
}

func Benchmark_MergeLocations(b *testing.B) {
	const locationCount = 1000000
	left := make([]*model.Location, locationCount+1)
	right := make([]*model.Location, locationCount+1)

	// Duplicates
	for i := 1; i < locationCount+1; i++ {
		left[i] = &model.Location{
			LocationID:     i,
			BookNumber:     sql.NullInt32{Int32: int32(i), Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: int32(i), Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: fmt.Sprint(i), Valid: true},
		}
		right[i] = &model.Location{
			LocationID:     i,
			BookNumber:     sql.NullInt32{Int32: int32(i), Valid: true},
			ChapterNumber:  sql.NullInt32{Int32: int32(i), Valid: true},
			DocumentID:     sql.NullInt32{},
			Track:          sql.NullInt32{},
			IssueTagNumber: 0,
			KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
			MepsLanguage:   2,
			LocationType:   0,
			Title:          sql.NullString{String: fmt.Sprint(i), Valid: true},
		}
	}

	MergeLocations(left, right)
}

func Test_solveLocationMergeConflict(t *testing.T) {
	conflicts := map[string]mergeConflict{
		"ChooseLeftConflict": {
			left: &model.Location{
				LocationID:     1,
				BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseLeft", Valid: true},
			},
			right: &model.Location{
				LocationID:     5,
				BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
		},
		"ChooseRightConflict": {
			left: &model.Location{
				LocationID:     2,
				BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 2, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
			right: &model.Location{
				LocationID:     6,
				BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 2, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseRight", Valid: true},
			},
		},
		"ChooseRightBecauseBothEmpty": {
			left: &model.Location{
				LocationID:     3,
				BookNumber:     sql.NullInt32{Int32: 3, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
			right: &model.Location{
				LocationID:     7,
				BookNumber:     sql.NullInt32{Int32: 3, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
		},
	}

	expectedResult := map[string]mergeSolution{
		"ChooseLeftConflict": {
			side: leftSide,
			solution: &model.Location{
				LocationID:     1,
				BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseLeft", Valid: true},
			},
			discarded: &model.Location{
				LocationID:     5,
				BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
		},
		"ChooseRightConflict": {
			side: rightSide,
			solution: &model.Location{
				LocationID:     6,
				BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 2, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseRight", Valid: true},
			},
			discarded: &model.Location{
				LocationID:     2,
				BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 2, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
		},
		"ChooseRightBecauseBothEmpty": {
			side: rightSide,
			solution: &model.Location{
				LocationID:     7,
				BookNumber:     sql.NullInt32{Int32: 3, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
			discarded: &model.Location{
				LocationID:     3,
				BookNumber:     sql.NullInt32{Int32: 3, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "", Valid: true},
			},
		},
	}

	result, _ := solveLocationMergeConflict(conflicts)

	assert.Equal(t, expectedResult, result)

	assert.PanicsWithValue(t, "No other type than *model.Location is supported! Given: *model.Bookmark", func() {
		panicConflict := map[string]MergeConflict{
			"WrongType": {
				left:  &model.Bookmark{},
				right: &model.Bookmark{},
			},
		}
		solveLocationMergeConflict(panicConflict)
	})
}
