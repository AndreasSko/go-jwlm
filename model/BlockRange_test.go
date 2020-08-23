package model

import (
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
