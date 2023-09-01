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

	assert.Equal(t, "VERYUNIQUEID_1_1_1_2_1_1_20_15_25_1", m1.UniqueKey())
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

func TestUserMarkBlockRange_PrettyPrint(t *testing.T) {
	m1 := &UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   1,
			ColorIndex:   5,
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
				EndToken:     sql.NullInt32{4, true},
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

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, "\n\nColorIndex:\t5\n")
	fmt.Fprint(w, "\nIdentifier:\t1\n")
	fmt.Fprint(w, "StartToken:\t0\n")
	fmt.Fprint(w, "EndToken:\t5\n\n")
	fmt.Fprint(w, "Identifier:\t2\n")
	fmt.Fprint(w, "StartToken:\t0\n")
	fmt.Fprint(w, "EndToken:\t4\n\n")
	fmt.Fprint(w, "Identifier:\t3\n")
	fmt.Fprint(w, "StartToken:\t0\n")
	fmt.Fprint(w, "EndToken:\t20\n")
	w.Flush()
	expectedResult := buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))

	db := &Database{
		Location: []*Location{
			nil,
			{
				LocationID:   1,
				Title:        sql.NullString{"Location-Title", true},
				MepsLanguage: sql.NullInt32{Int32: 0, Valid: true},
			},
		},
	}

	buf.Reset()
	fmt.Fprint(w, "\nTitle:\tLocation-Title\nIssueTagNumber:\t0\nMepsLanguage:\t0")
	fmt.Fprint(w, "\n\nColorIndex:\t5\n")
	fmt.Fprint(w, "\nIdentifier:\t1\n")
	fmt.Fprint(w, "StartToken:\t0\n")
	fmt.Fprint(w, "EndToken:\t5\n\n")
	fmt.Fprint(w, "Identifier:\t2\n")
	fmt.Fprint(w, "StartToken:\t0\n")
	fmt.Fprint(w, "EndToken:\t4\n\n")
	fmt.Fprint(w, "Identifier:\t3\n")
	fmt.Fprint(w, "StartToken:\t0\n")
	fmt.Fprint(w, "EndToken:\t20\n")
	w.Flush()
	expectedResult = buf.String()
	assert.Equal(t, expectedResult, m1.PrettyPrint(db))
}

func TestUserMarkBlockRange_RelatedEntries(t *testing.T) {
	db := &Database{
		Location: []*Location{
			nil,
			{
				LocationID: 1,
				Title:      sql.NullString{"Location-Title", true},
			},
		},
	}
	m1 := &UserMarkBlockRange{
		UserMark: &UserMark{
			UserMarkID:   1,
			ColorIndex:   5,
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
				EndToken:     sql.NullInt32{4, true},
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

	assert.Equal(t, Related{}, m1.RelatedEntries(nil))
	assert.Equal(t, Related{Location: db.Location[1]}, m1.RelatedEntries(db))
}

func TestUserMarkBlockRange_MarshalJSON(t *testing.T) {
	m1 := &UserMarkBlockRange{
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

	result, err := json.Marshal(m1)
	assert.NoError(t, err)
	assert.Equal(t,
		`{"type":"UserMarkBlockRange","userMark":{"type":"UserMark","userMarkId":12345,"colorIndex":0,"locationId":0,"styleIndex":0,"userMarkGuid":"VERYUNIQUEID","version":0},"blockRanges":[{"type":"BlockRange","blockRangeId":1,"blockType":1,"identifier":1,"startToken":{"Int32":1,"Valid":true},"endToken":{"Int32":2,"Valid":true},"userMarkId":1},{"type":"BlockRange","blockRangeId":2,"blockType":1,"identifier":20,"startToken":{"Int32":15,"Valid":true},"endToken":{"Int32":25,"Valid":true},"userMarkId":1}]}`,
		string(result))
}
