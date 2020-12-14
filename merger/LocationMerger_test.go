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
	// Check if original has not been tweaked
	assert.Equal(t, 6, left[6].LocationID)
	assert.Equal(t, 7, right[7].LocationID)
}

func Test_MergeLocationsDiffBibleEditions(t *testing.T) {
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
			KeySymbol:      sql.NullString{String: "nwt", Valid: true},
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
			KeySymbol:      sql.NullString{String: "nwt", Valid: true},
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
			KeySymbol:      sql.NullString{String: "nwt", Valid: true},
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
			KeySymbol:      sql.NullString{String: "nwt", Valid: true},
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

	result, _, _ := MergeLocations(left, right)

	assert.Equal(t, expectedResult, result)
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
	conflicts := map[string]MergeConflict{
		"ChooseLeftConflict": {
			Left: &model.Location{
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
			Right: &model.Location{
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
			Left: &model.Location{
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
			Right: &model.Location{
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
			Left: &model.Location{
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
			Right: &model.Location{
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

	expectedResult := map[string]MergeSolution{
		"ChooseLeftConflict": {
			Side: LeftSide,
			Solution: &model.Location{
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
			Discarded: &model.Location{
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
			Side: RightSide,
			Solution: &model.Location{
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
			Discarded: &model.Location{
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
			Side: RightSide,
			Solution: &model.Location{
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
			Discarded: &model.Location{
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
				Left:  &model.Bookmark{},
				Right: &model.Bookmark{},
			},
		}
		solveLocationMergeConflict(panicConflict)
	})
}

func Test_needsNwtstyMigration(t *testing.T) {
	type args struct {
		left  []*model.Location
		right []*model.Location
	}
	tests := []struct {
		args args
		want map[int]MergeSide
	}{
		{
			args: args{
				left: []*model.Location{
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"other", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"otherother", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 1},
				},
				right: []*model.Location{
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"other", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
				},
			},
			want: map[int]MergeSide{},
		},
		{
			args: args{
				left: []*model.Location{
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
				},
				right: []*model.Location{
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
				},
			},
			want: map[int]MergeSide{
				0: RightSide,
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, needsNwtstyMigration(tt.args.left, tt.args.right), tt.args)
	}
}

func Test_decideMigration(t *testing.T) {
	type args struct {
		langCounts map[int]*bibleEditionCounts
	}
	tests := []struct {
		args args
		want map[int]MergeSide
	}{
		{
			args: args{
				langCounts: map[int]*bibleEditionCounts{
					1: {
						leftNwt:     287,
						leftNwtsty:  2,
						rightNwt:    2,
						rightNwtsty: 34,
					},
				},
			},
			want: map[int]MergeSide{
				1: LeftSide,
			},
		},
		{
			args: args{
				langCounts: map[int]*bibleEditionCounts{
					0: {
						leftNwt:     1,
						leftNwtsty:  36,
						rightNwt:    0,
						rightNwtsty: 152,
					},
					1: {
						leftNwt:     1,
						leftNwtsty:  541,
						rightNwt:    2,
						rightNwtsty: 125,
					},
				},
			},
			want: map[int]MergeSide{},
		},
		{
			args: args{
				langCounts: map[int]*bibleEditionCounts{
					0: {
						leftNwt:     36,
						leftNwtsty:  1,
						rightNwt:    152,
						rightNwtsty: 0,
					},
					1: {
						leftNwt:     541,
						leftNwtsty:  1,
						rightNwt:    125,
						rightNwtsty: 2,
					},
				},
			},
			want: map[int]MergeSide{},
		},
		{
			args: args{
				langCounts: map[int]*bibleEditionCounts{
					0: {
						leftNwt:     1,
						leftNwtsty:  36,
						rightNwt:    152,
						rightNwtsty: 0,
					},
					1: {
						leftNwt:     541,
						leftNwtsty:  1,
						rightNwt:    125,
						rightNwtsty: 2,
					},
				},
			},
			want: map[int]MergeSide{
				0: RightSide,
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, decideMigration(tt.args.langCounts), tt.args)
	}
}

func Test_moveToNwtsty(t *testing.T) {
	type args struct {
		langs map[int]MergeSide
		left  []*model.Location
		right []*model.Location
	}
	tests := []struct {
		args args
		want args
	}{
		{
			args: args{
				langs: map[int]MergeSide{
					0: LeftSide,
					1: RightSide,
				},
				left: []*model.Location{
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"other", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"other", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 0},
				},
				right: []*model.Location{
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
				},
			},
			want: args{
				left: []*model.Location{
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"other", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"other", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
				},
				right: []*model.Location{
					{KeySymbol: sql.NullString{"nwt", true}, MepsLanguage: 0},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 1},
					{KeySymbol: sql.NullString{"nwtsty", true}, MepsLanguage: 0},
				},
			},
		},
	}
	for _, tt := range tests {
		moveToNwtsty(tt.args.langs, tt.args.left, tt.args.right)
		assert.Equal(t, tt.want.left, tt.args.left, tt.args.left)
		assert.Equal(t, tt.want.right, tt.args.right, tt.args.right)
	}
}

func Test_significantlyHigher(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		args args
		want bool
	}{
		{
			args: args{a: 1, b: 2},
			want: false,
		},
		{
			args: args{a: 10, b: 2},
			want: false,
		},
		{
			args: args{a: 25, b: 2},
			want: true,
		},
		{
			args: args{a: 11, b: 1},
			want: true,
		},
		{
			args: args{a: 1, b: 20},
			want: false,
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, significantlyHigher(tt.args.a, tt.args.b), tt.args)
	}
}
