package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"text/tabwriter"

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

func TestNote_UniqueKey(t *testing.T) {
	m := &Note{NoteID: 1, GUID: "A GUID"}
	assert.Equal(t, "A GUID", m.UniqueKey())
}

func TestNote_Equals(t *testing.T) {
	m1 := &Note{
		NoteID:          1,
		GUID:            "GUIDFOR1",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A Title", Valid: true},
		Content:         sql.NullString{String: "The content", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	m1_1 := &Note{
		NoteID:          2,
		GUID:            "GUIDFOR1",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A Title", Valid: true},
		Content:         sql.NullString{String: "The content", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	m2 := &Note{
		NoteID:          3,
		GUID:            "GUIDFOR3",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A early Title", Valid: true},
		Content:         sql.NullString{String: "The early content", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	m2_1 := &Note{
		NoteID:          3,
		GUID:            "GUIDFOR3",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A early Title", Valid: true},
		Content:         sql.NullString{String: "The early content", Valid: true},
		LastModified:    "2020-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}
	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m2_1))
}

func TestNote_PrettyPrint(t *testing.T) {
	m1 := &Note{
		NoteID:          1,
		GUID:            "GUIDFOR1",
		UserMarkID:      sql.NullInt32{Int32: 1, Valid: true},
		LocationID:      sql.NullInt32{Int32: 1, Valid: true},
		Title:           sql.NullString{String: "A Title", Valid: true},
		Content:         sql.NullString{String: "A very long content string that should hopefully result in a line break after max. 80 characters...", Valid: true},
		LastModified:    "2017-06-01T19:36:28+0200",
		BlockType:       0,
		BlockIdentifier: sql.NullInt32{},
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, "\nTitle:\tA Title")
	fmt.Fprint(w, "\nContent:\tA very long content string that should hopefully result in a line\n\tbreak after max. 80 characters...")
	fmt.Fprint(w, "\nLastModified:\t2017-06-01T19:36:28+0200")
	w.Flush()
	expectedResult := buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))

	db := &Database{
		Location: []*Location{
			nil,
			{
				LocationID: 1,
				Title:      sql.NullString{"Location-Title", true},
			},
		},
		UserMark: []*UserMark{
			nil,
			{
				UserMarkID: 1,
				ColorIndex: 5,
			},
		},
	}

	buf.Reset()
	fmt.Fprint(w, "\nTitle:\tA Title")
	fmt.Fprint(w, "\nContent:\tA very long content string that should hopefully result in a line\n\tbreak after max. 80 characters...")
	fmt.Fprint(w, "\nLastModified:\t2017-06-01T19:36:28+0200")
	fmt.Fprint(w, "\n\n\nRelated Location:\n\nTitle:\tLocation-Title\nIssueTagNumber:\t0\nMepsLanguage:\t0")
	fmt.Fprint(w, "\n\n\nRelated UserMark:\n\nColorIndex:\t5")
	w.Flush()
	expectedResult = buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(db))
}
