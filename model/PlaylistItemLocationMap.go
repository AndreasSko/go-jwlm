package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// PlaylistItemLocationMap represents the PlaylistItemLocationMap table inside the JW Library database
type PlaylistItemLocationMap struct {
	PlaylistItemID      int
	LocationID          int
	MajorMultimediaType int
	BaseDurationTicks   sql.NullInt32
}

// ID returns the ID of the entry (composite primary key, returns 0)
func (m *PlaylistItemLocationMap) ID() int {
	return m.PlaylistItemID
}

// SetID sets the ID of the entry (no-op for composite primary key)
func (m *PlaylistItemLocationMap) SetID(id int) {
	// TODO: Do we need to set the ID here
}

// UniqueKey returns the key that makes this PlaylistItemLocationMap unique,
// so it can be used as a key in a map.
func (m *PlaylistItemLocationMap) UniqueKey() string {
	var sb strings.Builder
	sb.Grow(10)
	sb.WriteString(strconv.FormatInt(int64(m.PlaylistItemID), 10))
	sb.WriteString("_")
	sb.WriteString(strconv.FormatInt(int64(m.LocationID), 10))
	return sb.String()
}

// Equals checks if the PlaylistItemLocationMap is equal to the given one.
func (m *PlaylistItemLocationMap) Equals(m2 Model) bool {
	if m2, ok := m2.(*PlaylistItemLocationMap); ok {
		return m.PlaylistItemID == m2.PlaylistItemID &&
			m.LocationID == m2.LocationID &&
			m.MajorMultimediaType == m2.MajorMultimediaType &&
			m.BaseDurationTicks == m2.BaseDurationTicks
	}
	return false
}

// RelatedEntries returns entries that are related to this PlaylistItemLocationMap
func (m *PlaylistItemLocationMap) RelatedEntries(db *Database) Related {
	result := Related{}
	return result
}

// PrettyPrint returns a string representation mainly for debugging purposes
func (m *PlaylistItemLocationMap) PrettyPrint(db *Database) string {
	return fmt.Sprintf("PlaylistItemLocationMap: PlaylistItem %d -> Location %d", m.PlaylistItemID, m.LocationID)
}

// tableName returns the name of the table
func (m *PlaylistItemLocationMap) tableName() string {
	return "PlaylistItemLocationMap"
}

// idName returns the name of the ID field
func (m *PlaylistItemLocationMap) idName() string {
	return "PlaylistItemId"
}

// scanRow scans a database row into the PlaylistItemLocationMap struct
func (m *PlaylistItemLocationMap) scanRow(rows *sql.Rows) (Model, error) {
	var result PlaylistItemLocationMap
	err := rows.Scan(
		&result.PlaylistItemID,
		&result.LocationID,
		&result.MajorMultimediaType,
		&result.BaseDurationTicks,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarshalJSON returns the JSON encoding
func (m *PlaylistItemLocationMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		PlaylistItemID      int           `json:"PlaylistItemId"`
		LocationID          int           `json:"LocationId"`
		MajorMultimediaType int           `json:"MajorMultimediaType"`
		BaseDurationTicks   sql.NullInt32 `json:"BaseDurationTicks"`
	}{
		PlaylistItemID:      m.PlaylistItemID,
		LocationID:          m.LocationID,
		MajorMultimediaType: m.MajorMultimediaType,
		BaseDurationTicks:   m.BaseDurationTicks,
	})
}

// UnmarshalJSON parses the JSON encoding
func (m *PlaylistItemLocationMap) UnmarshalJSON(data []byte) error {
	aux := &struct {
		PlaylistItemID      int           `json:"PlaylistItemId"`
		LocationID          int           `json:"LocationId"`
		MajorMultimediaType int           `json:"MajorMultimediaType"`
		BaseDurationTicks   sql.NullInt32 `json:"BaseDurationTicks"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	m.PlaylistItemID = aux.PlaylistItemID
	m.LocationID = aux.LocationID
	m.MajorMultimediaType = aux.MajorMultimediaType
	m.BaseDurationTicks = aux.BaseDurationTicks
	return nil
}

// MakeSlice converts a slice of the generic interface model
func (PlaylistItemLocationMap) MakeSlice(mdl []Model) []*PlaylistItemLocationMap {
	result := make([]*PlaylistItemLocationMap, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(*PlaylistItemLocationMap)
		}
	}
	return result
}
