package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockRange_SetID(t *testing.T) {
	m1 := &BlockRange{BlockRangeID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.BlockRangeID)

	m2 := BlockRange{BlockRangeID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.BlockRangeID)
}

func TestBlockRange_UniqueKey(t *testing.T) {
	m1 := &BlockRange{
		BlockRangeID: 1,
		BlockType:    1,
		Identifier:   1,
		StartToken:   sql.NullInt32{Int32: 1, Valid: true},
		EndToken:     sql.NullInt32{Int32: 2, Valid: true},
		UserMarkID:   1,
	}
	assert.Equal(t, "1_1_1_2", m1.UniqueKey())

	m2 := &BlockRange{
		BlockRangeID: 2,
		BlockType:    1,
		Identifier:   20,
		StartToken:   sql.NullInt32{Int32: 15, Valid: true},
		EndToken:     sql.NullInt32{Int32: 25, Valid: true},
		UserMarkID:   3334,
	}
	assert.Equal(t, "1_20_15_25", m2.UniqueKey())
}
