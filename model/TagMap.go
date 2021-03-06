package model

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
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
	var sb strings.Builder
	sb.Grow(15)
	sb.WriteString(strconv.FormatInt(int64(m.PlaylistItemID.Int32), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.LocationID.Int32), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.NoteID.Int32), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.TagID), 10))
	return sb.String()
}

// Equals checks if the TagMap is equal to the given one.
func (m *TagMap) Equals(m2 Model) bool {
	if m2, ok := m2.(*TagMap); ok {
		return m.PlaylistItemID.Int32 == m2.PlaylistItemID.Int32 &&
			m.LocationID.Int32 == m2.LocationID.Int32 &&
			m.NoteID.Int32 == m2.NoteID.Int32 &&
			m.TagID == m2.TagID &&
			m.Position == m2.Position
	}

	return false
}

// RelatedEntries returns entries that are related to this one
func (m *TagMap) RelatedEntries(db *Database) Related {
	// We don't need it for now, so just return empty slice
	return Related{}
}

// PrettyPrint prints TagMap in a human readable format and
// adds information about related entries if helpful.
func (m *TagMap) PrettyPrint(db *Database) string {
	panic("Not supported")
}

// MarshalJSON returns the JSON encoding of the entry
func (m TagMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type           string        `json:"type"`
		TagMapID       int           `json:"tagMapId"`
		PlaylistItemID sql.NullInt32 `json:"playlistItemId"`
		LocationID     sql.NullInt32 `json:"locationId"`
		NoteID         sql.NullInt32 `json:"noteId"`
		TagID          int           `json:"tagId"`
		Position       int           `json:"position"`
	}{
		Type:           "TagMap",
		TagMapID:       m.TagMapID,
		PlaylistItemID: m.PlaylistItemID,
		LocationID:     m.LocationID,
		NoteID:         m.NoteID,
		TagID:          m.TagID,
		Position:       m.Position,
	})
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
