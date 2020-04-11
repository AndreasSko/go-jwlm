package model

import "database/sql"

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
func (m BlockRange) ID() int {
	return m.BlockRangeID
}

func (m BlockRange) tableName() string {
	return "BlockRange"
}

func (m BlockRange) idName() string {
	return "BlockRangeId"
}

func (m BlockRange) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.BlockRangeID, &m.BlockType, &m.Identifier, &m.StartToken, &m.EndToken, &m.UserMarkID)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (BlockRange) makeSlice(mdl []model) []BlockRange {
	result := make([]BlockRange, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(BlockRange)
		}
	}
	return result
}
