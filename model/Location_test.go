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

func TestLocation_SetID(t *testing.T) {
	m1 := &Location{LocationID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.LocationID)

	m2 := Location{LocationID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.LocationID)
}

func TestLocation_UniqueKey(t *testing.T) {
	m1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "ThisTitleShouldNotBeInUniqueKey", Valid: true},
	}

	m2 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{},
		ChapterNumber:  sql.NullInt32{},
		DocumentID:     sql.NullInt32{},
		Track:          sql.NullInt32{},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "ThisOTitleShouldNotBeInUniqueKeyEither", Valid: true},
	}

	m3 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{},
		ChapterNumber:  sql.NullInt32{},
		DocumentID:     sql.NullInt32{},
		Track:          sql.NullInt32{},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{},
		MepsLanguage:   sql.NullInt32{},
		LocationType:   8,
		Title:          sql.NullString{String: "ThisOTitleShouldNotBeInUniqueKeyEither", Valid: true},
	}

	assert.Equal(t, "2_3_4_5_6_nwtsty_7_8", m1.UniqueKey())
	assert.Equal(t, "0_0_0_0_6__7_8", m2.UniqueKey())
	assert.Equal(t, "0_0_0_0_6__!_8", m3.UniqueKey())
}

func TestLocation_PrettyPrint(t *testing.T) {
	m1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "A title", Valid: true},
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, "\nTitle:\tA title")
	fmt.Fprint(w, "\nBookNumber:\t2")
	fmt.Fprint(w, "\nChapterNumber:\t3")
	fmt.Fprint(w, "\nDocumentID:\t4")
	fmt.Fprint(w, "\nTrack:\t5")
	fmt.Fprint(w, "\nIssueTagNumber:\t6")
	fmt.Fprint(w, "\nKeySymbol:\tnwtsty")
	fmt.Fprint(w, "\nMepsLanguage:\t7")
	w.Flush()
	expectedResult := buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))

	m1.KeySymbol.Valid = false

	buf.Reset()
	fmt.Fprint(w, "\nTitle:\tA title")
	fmt.Fprint(w, "\nBookNumber:\t2")
	fmt.Fprint(w, "\nChapterNumber:\t3")
	fmt.Fprint(w, "\nDocumentID:\t4")
	fmt.Fprint(w, "\nTrack:\t5")
	fmt.Fprint(w, "\nIssueTagNumber:\t6")
	fmt.Fprint(w, "\nMepsLanguage:\t7")
	w.Flush()
	expectedResult = buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))
}

func TestLocation_Equals(t *testing.T) {
	m1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "ThisTitleShouldNotBeInUniqueKey", Valid: true},
	}
	m1_1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
	}

	m2 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{},
		ChapterNumber:  sql.NullInt32{},
		DocumentID:     sql.NullInt32{},
		Track:          sql.NullInt32{},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "ThisOTitleShouldNotBeInUniqueKeyEither", Valid: true},
	}

	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
}

func TestLocation_RelatedEntries(t *testing.T) {
	m1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "A title", Valid: true},
	}

	assert.Equal(t, Related{}, m1.RelatedEntries(nil))
	assert.Equal(t, Related{}, m1.RelatedEntries(&Database{}))
}

func TestLocation_MarshalJSON(t *testing.T) {
	m1 := &Location{
		LocationID:     1,
		BookNumber:     sql.NullInt32{Int32: 2, Valid: true},
		ChapterNumber:  sql.NullInt32{Int32: 3, Valid: true},
		DocumentID:     sql.NullInt32{Int32: 4, Valid: true},
		Track:          sql.NullInt32{Int32: 5, Valid: true},
		IssueTagNumber: 6,
		KeySymbol:      sql.NullString{String: "nwtsty", Valid: true},
		MepsLanguage:   sql.NullInt32{Int32: 7, Valid: true},
		LocationType:   8,
		Title:          sql.NullString{String: "ThisTitleShouldNotBeInUniqueKey", Valid: true},
	}

	result, err := json.Marshal(m1)
	assert.NoError(t, err)
	assert.Equal(t,
		`{"type":"Location","locationId":1,"bookNumber":{"Int32":2,"Valid":true},"chapterNumber":{"Int32":3,"Valid":true},"documentId":{"Int32":4,"Valid":true},"track":{"Int32":5,"Valid":true},"issueTagNumber":6,"keySymbol":{"String":"nwtsty","Valid":true},"mepsLanguage":{"Int32":7,"Valid":true},"locationType":8,"title":{"String":"ThisTitleShouldNotBeInUniqueKey","Valid":true}}`,
		string(result))
}
