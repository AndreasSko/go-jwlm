package model

import "database/sql"

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

func (m *Note) tableName() string {
	return "Note"
}

func (m *Note) idName() string {
	return "NoteId"
}

func (m *Note) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.NoteID, &m.GUID, &m.UserMarkID, &m.LocationID, &m.Title, &m.Content,
		&m.LastModified, &m.BlockType, &m.BlockIdentifier)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (Note) makeSlice(mdl []*model) []*Note {
	result := make([]*Note, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = (*mdl[i]).(*Note)
		}
	}
	return result
}
