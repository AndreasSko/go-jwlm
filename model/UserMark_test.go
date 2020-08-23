package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserMark_SetID(t *testing.T) {
	m1 := &UserMark{UserMarkID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.UserMarkID)

	m2 := UserMark{UserMarkID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.UserMarkID)
}
