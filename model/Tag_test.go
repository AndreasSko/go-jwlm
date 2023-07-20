package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"text/tabwriter"

	"github.com/stretchr/testify/assert"
)

func TestTag_SetID(t *testing.T) {
	m1 := &Tag{TagID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.TagID)

	m2 := Tag{TagID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.TagID)
}

func TestTag_UniqueKey(t *testing.T) {
	m1 := &Tag{
		TagID:   1,
		TagType: 1,
		Name:    "FirstTag",
	}
	assert.Equal(t, "1_FirstTag", m1.UniqueKey())

	m2 := &Tag{
		TagID:   1,
		TagType: 2000000000,
		Name:    "Another Tag with spaces",
	}
	assert.Equal(t, "2000000000_Another Tag with spaces", m2.UniqueKey())
}

func TestTag_Equals(t *testing.T) {
	m1 := &Tag{
		TagID:   1,
		TagType: 1,
		Name:    "FirstTag",
	}
	m1_1 := &Tag{
		TagID:   100000,
		TagType: 1,
		Name:    "FirstTag",
	}
	m2 := &Tag{
		TagID:   2,
		TagType: 1,
		Name:    "Different Tag",
	}
	m2_1 := &Tag{
		TagID:   200000,
		TagType: 1,
		Name:    "Different Tag",
	}
	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m2_1))

}

func TestTag_PrettyPrint(t *testing.T) {
	m1 := &Tag{
		TagID:   1,
		TagType: 1,
		Name:    "First  Tag",
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, "\nName:\tFirst  Tag")
	w.Flush()
	expectedResult := buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))
}

func TestTag_RelatedEntries(t *testing.T) {
	m1 := &Tag{
		TagID:   1,
		TagType: 1,
		Name:    "FirstTag",
	}

	assert.Equal(t, Related{}, m1.RelatedEntries(nil))
	assert.Equal(t, Related{}, m1.RelatedEntries(&Database{}))
}

func TestTag_MarshalJSON(t *testing.T) {
	m1 := &Tag{
		TagID:   1,
		TagType: 2,
		Name:    "FirstTag",
	}

	result, err := json.Marshal(m1)
	assert.NoError(t, err)
	assert.Equal(t,
		`{"type":"Tag","tagId":1,"tagType":2,"name":"FirstTag"}`,
		string(result))
}
