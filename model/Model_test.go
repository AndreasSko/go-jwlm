package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyPrint(t *testing.T) {
	location := &Location{
		LocationID: 1,
	}
	assert.PanicsWithValue(t, "Given struct does not contain field notexistent", func() {
		prettyPrint(location, []string{"notexistent"})
	})

	umbr := &UserMarkBlockRange{
		UserMark: &UserMark{},
	}

	assert.PanicsWithValue(t, "Unsupported type for field UserMark", func() {
		prettyPrint(umbr, []string{"UserMark"})
	})
}
