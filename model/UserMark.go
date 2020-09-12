package model

import (
	"database/sql"
)

// UserMark represents the UserMark table inside the JW Library database
type UserMark struct {
	UserMarkID   int
	ColorIndex   int
	LocationID   int
	StyleIndex   int
	UserMarkGUID string
	Version      int
}

// ID returns the ID of the entry
func (m *UserMark) ID() int {
	return m.UserMarkID
}

// SetID sets the ID of the entry
func (m *UserMark) SetID(id int) {
	m.UserMarkID = id
}

// UniqueKey returns the key that makes this UserMark unique,
// so it can be used as a key in a map.
func (m *UserMark) UniqueKey() string {
	return m.UserMarkGUID
}

// Equals checks if the UserMark is equal to the given one.
func (m *UserMark) Equals(m2 Model) bool {
	if m2, ok := m2.(*UserMark); ok {
		return m.ColorIndex == m2.ColorIndex &&
			m.LocationID == m2.LocationID &&
			m.StyleIndex == m2.StyleIndex &&
			m.Version == m2.Version
	}
	return false
}

func (m *UserMark) tableName() string {
	return "UserMark"
}

func (m *UserMark) idName() string {
	return "UserMarkId"
}

func (m *UserMark) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.UserMarkID, &m.ColorIndex, &m.LocationID, &m.StyleIndex, &m.UserMarkGUID, &m.Version)
	return m, err
}

// MakeSlice converts a slice of the generice interface model
func (UserMark) MakeSlice(mdl []Model) []*UserMark {
	result := make([]*UserMark, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(*UserMark)
		}
	}
	return result
}
