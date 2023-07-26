package model

import (
	"database/sql"
	"encoding/json"
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
	Created         string
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
	return m.GUID
}

// Equals checks if the Note is equal to the given one.
func (m *Note) Equals(m2 Model) bool {
	if m2, ok := m2.(*Note); ok {
		return m.GUID == m2.GUID &&
			m.Title == m2.Title &&
			m.Content == m2.Content
	}
	return false
}

// RelatedEntries returns entries that are related to this one
func (m *Note) RelatedEntries(db *Database) Related {
	result := Related{}

	if location := db.FetchFromTable("Location", int(m.LocationID.Int32)); location != nil {
		result.Location = location.(*Location)
	}

	// Todo: Maybe add BlockRange or rather use UserMarkBlockRange?
	if userMark := db.FetchFromTable("UserMark", int(m.UserMarkID.Int32)); userMark != nil {
		result.UserMark = userMark.(*UserMark)
	}

	return result
}

// PrettyPrint prints Note in a human readable format and
// adds information about related entries if helpful.
func (m *Note) PrettyPrint(db *Database) string {
	fields := []string{"Title", "Content", "LastModified", "Created"}
	result := prettyPrint(m, fields)

	// TODO: Use RelatedEntries
	if location := db.FetchFromTable("Location", int(m.LocationID.Int32)); location != nil {
		result += "\n\n\nRelated Location:\n"
		result += location.PrettyPrint(db)
	}

	if userMark := db.FetchFromTable("UserMark", int(m.UserMarkID.Int32)); userMark != nil {
		result += "\n\n\nRelated UserMark:\n"
		result += userMark.PrettyPrint(db)
	}

	return result
}

// MarshalJSON returns the JSON encoding of the entry
func (m Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type            string         `json:"type"`
		NoteID          int            `json:"noteId"`
		GUID            string         `json:"guid"`
		UserMarkID      sql.NullInt32  `json:"userMarkId"`
		LocationID      sql.NullInt32  `json:"locationId"`
		Title           sql.NullString `json:"title"`
		Content         sql.NullString `json:"content"`
		LastModified    string         `json:"lastModified"`
		Created         string         `json:"created"`
		BlockType       int            `json:"blockType"`
		BlockIdentifier sql.NullInt32  `json:"blockIdentifier"`
	}{
		Type:            "Note",
		NoteID:          m.NoteID,
		GUID:            m.GUID,
		UserMarkID:      m.UserMarkID,
		LocationID:      m.LocationID,
		Title:           m.Title,
		Content:         m.Content,
		LastModified:    m.LastModified,
		Created:         m.Created,
		BlockType:       m.BlockType,
		BlockIdentifier: m.BlockIdentifier,
	})
}

func (m *Note) tableName() string {
	return "Note"
}

func (m *Note) idName() string {
	return "NoteId"
}

func (m *Note) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.NoteID, &m.GUID, &m.UserMarkID, &m.LocationID, &m.Title, &m.Content,
		&m.LastModified, &m.Created, &m.BlockType, &m.BlockIdentifier)
	return m, err
}
