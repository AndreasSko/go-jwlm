package model

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
)

// UserMarkBlockRange represents a UserMark joined with its BlockRange entries.
// It does NOT represent an actual table in the JWLibrary backup file!
type UserMarkBlockRange struct {
	UserMark    *UserMark
	BlockRanges []*BlockRange
}

// ID returns the ID of the UserMark representing the whole UserMarkBlockRange{}
func (m *UserMarkBlockRange) ID() int {
	return m.UserMark.ID()
}

// SetID sets the ID of the entry, which is the UserMarkID representing
// the whole UserMarkBlockRange{}
func (m *UserMarkBlockRange) SetID(id int) {
	m.UserMark.SetID(id)
}

// UniqueKey returns the key that makes this UserMarkBlockRange unique,
// so it can be used as a key in a map.
func (m *UserMarkBlockRange) UniqueKey() string {
	var sb strings.Builder
	sb.Grow(70)
	sb.WriteString(m.UserMark.UniqueKey())
	sb.WriteString("_")
	for i, br := range m.BlockRanges {
		sb.WriteString(br.UniqueKey())
		if i+1 < len(m.BlockRanges) {
			sb.WriteString("_")
		}
	}

	return sb.String()
}

var rmUID = regexp.MustCompile(`^(\d*_\d*_\d*_\d*)(_\d*)`)

// Equals checks if the UserMarkBlockRange is equal to the given one.
// It will both check its UserMark and all BlockRanges.
func (m *UserMarkBlockRange) Equals(m2 Model) bool {
	// Remove UserMarkID from UniqueKey, as BlockRanges have already
	// been joined with UserMark.
	// Compare UniqueKeys of both BlockRanges to check if they are the same
	mBRKeys := make(map[string]bool, len(m.BlockRanges))
	m2BRKeys := make(map[string]bool, len(m2.(*UserMarkBlockRange).BlockRanges))
	for _, br := range m.BlockRanges {
		uq := rmUID.FindStringSubmatch(br.UniqueKey())[1]
		mBRKeys[uq] = true
	}
	for _, br := range m2.(*UserMarkBlockRange).BlockRanges {
		uq := rmUID.FindStringSubmatch(br.UniqueKey())[1]
		m2BRKeys[uq] = true
	}

	return m.UserMark.Equals(m2.(*UserMarkBlockRange).UserMark) &&
		reflect.DeepEqual(mBRKeys, m2BRKeys)
}

// RelatedEntries returns entries that are related to this one
func (m *UserMarkBlockRange) RelatedEntries(db *Database) Related {
	result := Related{}

	if location := db.FetchFromTable("Location", m.UserMark.LocationID); location != nil {
		result.Location = location.(*Location)
	}

	return result
}

// PrettyPrint prints UserMarkBlockRange in a human readable format and
// adds information about related entries if helpful.
func (m *UserMarkBlockRange) PrettyPrint(db *Database) string {
	umFields := []string{"ColorIndex"}
	brFields := []string{"Identifier", "StartToken", "EndToken"}

	var result string

	if location := db.FetchFromTable("Location", m.UserMark.LocationID); location != nil {
		result += location.PrettyPrint(db)
	}

	result += "\n" + prettyPrint(m.UserMark, umFields) + "\n"

	for _, br := range m.BlockRanges {
		result += prettyPrint(br, brFields) + "\n"
	}

	return result
}

// MarshalJSON returns the JSON encoding of the entry
func (m UserMarkBlockRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type        string        `json:"type"`
		UserMark    *UserMark     `json:"userMark"`
		BlockRanges []*BlockRange `json:"blockRanges"`
	}{
		Type:        "UserMarkBlockRange",
		UserMark:    m.UserMark,
		BlockRanges: m.BlockRanges,
	})
}

func (m *UserMarkBlockRange) tableName() string {
	panic("Not supported!")
}

func (m *UserMarkBlockRange) idName() string {
	panic("Not supported!")
}

func (m *UserMarkBlockRange) scanRow(rows *sql.Rows) (Model, error) {
	panic("Not supported!")
}

// MakeSlice converts a slice of the generice interface model
func (UserMarkBlockRange) MakeSlice(mdl []Model) []*UserMarkBlockRange {
	panic("Not supported!")
}
