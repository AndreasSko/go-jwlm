package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserMarkBlockRange_ID(t *testing.T) {
	m1 := UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   12345,
			UserMarkGUID: "VERYUNIQUEID",
		},
		BlockRanges: []*BlockRange{
			{
				BlockRangeID: 1,
				BlockType:    1,
				Identifier:   1,
				StartToken:   sql.NullInt32{Int32: 1, Valid: true},
				EndToken:     sql.NullInt32{Int32: 2, Valid: true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 2,
				BlockType:    1,
				Identifier:   20,
				StartToken:   sql.NullInt32{Int32: 15, Valid: true},
				EndToken:     sql.NullInt32{Int32: 25, Valid: true},
				UserMarkID:   1,
			},
		},
	}
	assert.Equal(t, 12345, m1.ID())
}

func TestUserMarkBlockRange_SetID(t *testing.T) {
	m1 := UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   12345,
			UserMarkGUID: "VERYUNIQUEID",
		},
		BlockRanges: []*BlockRange{
			{
				BlockRangeID: 1,
				BlockType:    1,
				Identifier:   1,
				StartToken:   sql.NullInt32{Int32: 1, Valid: true},
				EndToken:     sql.NullInt32{Int32: 2, Valid: true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 2,
				BlockType:    1,
				Identifier:   20,
				StartToken:   sql.NullInt32{Int32: 15, Valid: true},
				EndToken:     sql.NullInt32{Int32: 25, Valid: true},
				UserMarkID:   1,
			},
		},
	}
	assert.Equal(t, 12345, m1.ID())
	m1.SetID(6789)
	assert.Equal(t, m1.ID(), 6789)
}

func TestUserMarkBlockRange_UniqueKey(t *testing.T) {
	m1 := UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkGUID: "VERYUNIQUEID",
		},
		BlockRanges: []*BlockRange{
			{
				BlockRangeID: 1,
				BlockType:    1,
				Identifier:   1,
				StartToken:   sql.NullInt32{Int32: 1, Valid: true},
				EndToken:     sql.NullInt32{Int32: 2, Valid: true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 2,
				BlockType:    1,
				Identifier:   20,
				StartToken:   sql.NullInt32{Int32: 15, Valid: true},
				EndToken:     sql.NullInt32{Int32: 25, Valid: true},
				UserMarkID:   1,
			},
		},
	}

	assert.Equal(t, "VERYUNIQUEID_1_1_1_2_1_20_15_25", m1.UniqueKey())
}

func TestUserMarkBlockRange_Equals(t *testing.T) {
	m1 := &UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   1,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   1,
			UserMarkGUID: "FIRST",
			Version:      1,
		},
		BlockRanges: []*BlockRange{
			{
				BlockRangeID: 1,
				BlockType:    1,
				Identifier:   1,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{5, true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 2,
				BlockType:    1,
				Identifier:   2,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{5, true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 3,
				BlockType:    1,
				Identifier:   3,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{20, true},
				UserMarkID:   1,
			},
		},
	}
	m1_1 := &UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   1000,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   1,
			UserMarkGUID: "FIRSTT",
			Version:      1,
		},
		BlockRanges: []*BlockRange{
			{
				BlockRangeID: 6,
				BlockType:    1,
				Identifier:   3,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{20, true},
				UserMarkID:   1000,
			},
			{
				BlockRangeID: 5,
				BlockType:    1,
				Identifier:   2,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{5, true},
				UserMarkID:   1000,
			},
			{
				BlockRangeID: 4,
				BlockType:    1,
				Identifier:   1,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{5, true},
				UserMarkID:   1000,
			},
		},
	}

	m2 := &UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   1,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   1,
			UserMarkGUID: "FIRST",
			Version:      1,
		},
		BlockRanges: []*BlockRange{
			{
				BlockRangeID: 1,
				BlockType:    1,
				Identifier:   1,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{5, true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 2,
				BlockType:    1,
				Identifier:   2,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{5, true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 3,
				BlockType:    1,
				Identifier:   3,
				StartToken:   sql.NullInt32{0, true},
				EndToken:     sql.NullInt32{21, true},
				UserMarkID:   1,
			},
		},
	}

	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
}