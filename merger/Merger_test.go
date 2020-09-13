package merger

import (
	"database/sql"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/stretchr/testify/assert"
)

func Test_sortMergeSolution(t *testing.T) {
	solution := []mergeSolution{
		{
			side: rightSide,
			solution: &model.Location{
				LocationID:     3,
				BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseRight", Valid: true},
			},
		},
		{
			side: leftSide,
			solution: &model.Location{
				LocationID:     3,
				BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 2, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseLeft", Valid: true},
			},
		},
		{
			side: leftSide,
			solution: &model.Location{
				LocationID:     2,
				BookNumber:     sql.NullInt32{Int32: 3, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseLeft", Valid: true},
			},
		},
		{
			side: rightSide,
			solution: &model.Location{
				LocationID:     1,
				BookNumber:     sql.NullInt32{Int32: 4, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 4, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseRight", Valid: true},
			},
		},
	}

	expectedResult := []mergeSolution{
		{
			side: rightSide,
			solution: &model.Location{
				LocationID:     1,
				BookNumber:     sql.NullInt32{Int32: 4, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 4, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseRight", Valid: true},
			},
		},
		{
			side: leftSide,
			solution: &model.Location{
				LocationID:     2,
				BookNumber:     sql.NullInt32{Int32: 3, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseLeft", Valid: true},
			},
		},
		{
			side: leftSide,
			solution: &model.Location{
				LocationID:     3,
				BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 2, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseLeft", Valid: true},
			},
		},
		{
			side: rightSide,
			solution: &model.Location{
				LocationID:     3,
				BookNumber:     sql.NullInt32{Int32: 1, Valid: true},
				ChapterNumber:  sql.NullInt32{Int32: 1, Valid: true},
				DocumentID:     sql.NullInt32{},
				Track:          sql.NullInt32{},
				IssueTagNumber: 0,
				KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
				MepsLanguage:   2,
				LocationType:   0,
				Title:          sql.NullString{String: "ChooseRight", Valid: true},
			},
		},
	}

	sortMergeSolution(&solution)

	assert.Equal(t, expectedResult, solution)
}
