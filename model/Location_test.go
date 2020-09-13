package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocation_SetID(t *testing.T) {
	m1 := &Location{LocationID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.LocationID)

	m2 := Location{LocationID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.LocationID)
}

func TestLocation_UniqueKey(t *testing.T) {
	m1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   7,
		LocationType:   8,
		Title:          sql.NullString{String: "ThisTitleShouldNotBeInUniqueKey", Valid: true},
	}

	m2 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{},
		ChapterNumber:  sql.NullInt32{},
		DocumentID:     sql.NullInt32{},
		Track:          sql.NullInt32{},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{},
		MepsLanguage:   7,
		LocationType:   8,
		Title:          sql.NullString{String: "ThisOTitleShouldNotBeInUniqueKeyEither", Valid: true},
	}

	assert.Equal(t, "2_3_4_5_6_nwtsty_7_8", m1.UniqueKey())
	assert.Equal(t, "0_0_0_0_6__7_8", m2.UniqueKey())
}
