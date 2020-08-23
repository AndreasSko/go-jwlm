package model

import "database/sql"

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

func (m *UserMark) tableName() string {
	return "UserMark"
}

func (m *UserMark) idName() string {
	return "UserMarkId"
}

func (m *UserMark) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.UserMarkID, &m.ColorIndex, &m.LocationID, &m.StyleIndex, &m.UserMarkGUID, &m.Version)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (UserMark) makeSlice(mdl []*model) []*UserMark {
	result := make([]*UserMark, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = (*mdl[i]).(*UserMark)
		}
	}
	return result
}
