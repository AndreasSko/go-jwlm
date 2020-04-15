package model

import "database/sql"

// TagMap represents the TagMap table inside the JW Library database
type TagMap struct {
	TagMapID       int
	PlaylistItemID sql.NullInt32
	LocationID     sql.NullInt32
	NoteID         sql.NullInt32
	TagID          int
	Position       int
}

// ID returns the ID of the entry
func (m TagMap) ID() int {
	return m.TagMapID
}

func (m TagMap) tableName() string {
	return "TagMap"
}

func (m TagMap) idName() string {
	return "TagMapId"
}

func (m TagMap) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.TagMapID, &m.PlaylistItemID, &m.LocationID, &m.NoteID, &m.TagID, &m.Position)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (TagMap) makeSlice(mdl []model) []TagMap {
	result := make([]TagMap, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(TagMap)
		}
	}
	return result
}
