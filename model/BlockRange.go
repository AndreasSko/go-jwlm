package model

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
)

// BlockRange represents the BlockRange table inside the JW Library database
type BlockRange struct {
	BlockRangeID int
	BlockType    int
	Identifier   int
	StartToken   sql.NullInt32
	EndToken     sql.NullInt32
	UserMarkID   int
}

// ID returns the ID of the entry
func (m *BlockRange) ID() int {
	return m.BlockRangeID
}

// SetID sets the ID of the entry
func (m *BlockRange) SetID(id int) {
	m.BlockRangeID = id
}

// UniqueKey returns the key that makes this BlockRange unique,
// so it can be used as a key in a map.
func (m *BlockRange) UniqueKey() string {
	var sb strings.Builder
	sb.Grow(15)
	sb.WriteString(strconv.FormatInt(int64(m.BlockType), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.Identifier), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.StartToken.Int32), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.EndToken.Int32), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.UserMarkID), 10))
	return sb.String()
}

// Equals checks if the BlockRange is equal to the given one.
func (m *BlockRange) Equals(m2 Model) bool {
	if m2, ok := m2.(*BlockRange); ok {
		return m.BlockType == m2.BlockType &&
			m.Identifier == m2.Identifier &&
			m.StartToken.Int32 == m2.StartToken.Int32 &&
			m.EndToken.Int32 == m2.EndToken.Int32 &&
			m.UserMarkID == m2.UserMarkID
	}
	return false
}

// RelatedEntries returns entries that are related to this one
func (m *BlockRange) RelatedEntries(db *Database) []Model {
	// We don't need it for now, so just return empty slice
	return []Model{}
}

// PrettyPrint prints BlockRange in a human readable format and
// adds information about related entries if helpful.
func (m *BlockRange) PrettyPrint(db *Database) string {
	fields := []string{"Identifier", "StartToken", "EndToken"}
	return prettyPrint(m, fields)
}

// MarshalJSON returns the JSON encoding of the entry
func (m BlockRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type         string
		BlockRangeID int
		BlockType    int
		Identifier   int
		StartToken   sql.NullInt32
		EndToken     sql.NullInt32
		UserMarkID   int
	}{
		Type:         "BlockRange",
		BlockRangeID: m.BlockRangeID,
		BlockType:    m.BlockType,
		Identifier:   m.Identifier,
		StartToken:   m.StartToken,
		EndToken:     m.EndToken,
		UserMarkID:   m.UserMarkID,
	})
}

func (m *BlockRange) tableName() string {
	return "BlockRange"
}

func (m *BlockRange) idName() string {
	return "BlockRangeId"
}

func (m *BlockRange) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.BlockRangeID, &m.BlockType, &m.Identifier, &m.StartToken, &m.EndToken, &m.UserMarkID)
	return m, err
}

// MakeSlice converts a slice of the generice interface model
func (BlockRange) MakeSlice(mdl []Model) []*BlockRange {
	result := make([]*BlockRange, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(*BlockRange)
		}
	}
	return result
}
