package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"text/tabwriter"

	"github.com/stretchr/testify/assert"
)

func TestBookmark_SetID(t *testing.T) {
	m1 := &Bookmark{BookmarkID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.BookmarkID)

	m2 := Bookmark{BookmarkID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.BookmarkID)
}

func TestBookmark_UniqueKey(t *testing.T) {
	m1 := &Bookmark{
		BookmarkID:            1,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}

	assert.Equal(t, "3_4", m1.UniqueKey())
}

func TestBookmark_Equals(t *testing.T) {
	m1 := &Bookmark{
		BookmarkID:            1,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}
	m1_1 := &Bookmark{
		BookmarkID:            1000,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}
	m2 := &Bookmark{
		BookmarkID:            1000,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  5,
		Title:                 "Test",
		Snippet:               sql.NullString{},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}

	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
}

func TestBookmark_PrettyPrint(t *testing.T) {
	m1 := &Bookmark{
		BookmarkID:            1,
		LocationID:            1,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{"A snippet", true},
		BlockType:             0,
		BlockIdentifier:       sql.NullInt32{},
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, "\nTitle:\tTest")
	fmt.Fprint(w, "\nSnippet:\tA snippet")
	fmt.Fprint(w, "\nSlot:\t4")
	fmt.Fprint(w, "\nPublicationLocationID:\t3")
	w.Flush()
	expectedResult := buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))

	m1.Title = ""
	db := &Database{
		Location: []*Location{
			nil,
			{
				LocationID: 1,
				Title:      sql.NullString{"Location-Title", true},
			},
		},
	}

	buf.Reset()
	fmt.Fprint(w, "\nTitle:\t")
	fmt.Fprint(w, "\nSnippet:\tA snippet")
	fmt.Fprint(w, "\nSlot:\t4")
	fmt.Fprint(w, "\nPublicationLocationID:\t3")
	fmt.Fprint(w, "\n\n\nRelated Location:\n\nTitle:\tLocation-Title\nIssueTagNumber:\t0\nMepsLanguage:\t0")
	w.Flush()
	expectedResult = buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(db))
}

func TestBookmark_RelatedEntries(t *testing.T) {
	db := &Database{
		Bookmark: []*Bookmark{
			nil,
			{
				BookmarkID:            1,
				LocationID:            1,
				PublicationLocationID: 3,
				Slot:                  4,
				Title:                 "Test",
				Snippet:               sql.NullString{},
				BlockType:             0,
				BlockIdentifier:       sql.NullInt32{},
			},
		},
		Location: []*Location{
			nil,
			{
				LocationID: 1,
				Title:      sql.NullString{"Location-Title", true},
			},
		},
	}

	assert.Equal(t, Related{}, db.Bookmark[1].RelatedEntries(nil))
	assert.Equal(t,
		Related{Location: db.Location[1], PublicationLocation: db.Location[1]},
		db.Bookmark[1].RelatedEntries(db))
}

func TestBookmark_MarshalJSON(t *testing.T) {
	m1 := &Bookmark{
		BookmarkID:            1,
		LocationID:            2,
		PublicationLocationID: 3,
		Slot:                  4,
		Title:                 "Test",
		Snippet:               sql.NullString{"A snippet", true},
		BlockType:             5,
		BlockIdentifier:       sql.NullInt32{},
	}

	result, err := json.Marshal(m1)
	assert.NoError(t, err)
	assert.Equal(t,
		`{"Type":"Bookmark","BookmarkID":1,"LocationID":2,"PublicationLocationID":3,"Slot":4,"Title":"Test","Snippet":{"String":"A snippet","Valid":true},"BlockType":5,"BlockIdentifier":{"Int32":0,"Valid":false}}`,
		string(result))
}
