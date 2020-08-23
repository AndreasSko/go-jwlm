package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNote_SetID(t *testing.T) {
	m1 := &Note{NoteID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.NoteID)

	m2 := Note{NoteID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.NoteID)
}
