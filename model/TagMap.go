package model

import (
	"database/sql"
	"fmt"
)

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
func (m *TagMap) ID() int {
	return m.TagMapID
}

// SetID sets the ID of the entry
func (m *TagMap) SetID(id int) {
	m.TagMapID = id
}

// UniqueKey returns the key that makes this TagMap unique,
// so it can be used as a key in a map.
func (m *TagMap) UniqueKey() string {
	return fmt.Sprintf("TODO")
}

// Equals checks if the TagMap is equal to the given one.
func (m *TagMap) Equals(m2 Model) bool {
	return false
}

func (m *TagMap) tableName() string {
	return "TagMap"
}

func (m *TagMap) idName() string {
	return "TagMapId"
}

func (m *TagMap) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.TagMapID, &m.PlaylistItemID, &m.LocationID, &m.NoteID, &m.TagID, &m.Position)
	return m, err
}

// MakeSlice converts a slice of the generice interface model
func (TagMap) MakeSlice(mdl []Model) []*TagMap {
	result := make([]*TagMap, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(*TagMap)
		}
	}
	return result
}
