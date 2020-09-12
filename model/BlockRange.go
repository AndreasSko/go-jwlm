package model

import (
	"database/sql"
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
	return sb.String()
}

// Equals checks if the BlockRange is equal to the given one.
func (m *BlockRange) Equals(m2 Model) bool {
	return false
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
