package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagMap_SetID(t *testing.T) {
	m1 := &TagMap{TagMapID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.TagMapID)

	m2 := TagMap{TagMapID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.TagMapID)
}
