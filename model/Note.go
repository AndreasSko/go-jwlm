package model

import (
	"database/sql"
	"fmt"
)

// Note represents the Note table inside the JW Library database
type Note struct {
	NoteID          int
	GUID            string
	UserMarkID      sql.NullInt32
	LocationID      sql.NullInt32
	Title           sql.NullString
	Content         sql.NullString
	LastModified    string
	BlockType       int
	BlockIdentifier sql.NullInt32
}

// ID returns the ID of the entry
func (m *Note) ID() int {
	return m.NoteID
}

// SetID sets the ID of the entry
func (m *Note) SetID(id int) {
	m.NoteID = id
}

// UniqueKey returns the key that makes this Note unique,
// so it can be used as a key in a map.
func (m *Note) UniqueKey() string {
	return fmt.Sprintf("TODO")
}

// Equals checks if the Note is equal to the given one.
func (m *Note) Equals(m2 Model) bool {
	return false
}

func (m *Note) tableName() string {
	return "Note"
}

func (m *Note) idName() string {
	return "NoteId"
}

func (m *Note) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.NoteID, &m.GUID, &m.UserMarkID, &m.LocationID, &m.Title, &m.Content,
		&m.LastModified, &m.BlockType, &m.BlockIdentifier)
	return m, err
}

// MakeSlice converts a slice of the generice interface model
func (Note) MakeSlice(mdl []Model) []*Note {
	result := make([]*Note, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(*Note)
		}
	}
	return result
}
