package model

import (
	"database/sql"
)

// Location represents the Location table inside the JW Library database
type Location struct {
	LocationID     int
	BookNumber     sql.NullInt32
	ChapterNumber  sql.NullInt32
	DocumentID     sql.NullInt32
	Track          sql.NullInt32
	IssueTagNumber int
	KeySymbol      string
	MepsLanguage   int
	LocationType   int
	Title          sql.NullString
}

// ID returns the ID of the entry
func (m Location) ID() int {
	return m.LocationID
}

func (m Location) tableName() string {
	return "Location"
}

func (m Location) idName() string {
	return "LocationId"
}

func (m Location) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.LocationID, &m.BookNumber, &m.ChapterNumber, &m.DocumentID, &m.Track,
		&m.IssueTagNumber, &m.KeySymbol, &m.MepsLanguage, &m.LocationType, &m.Title)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (Location) makeSlice(mdl []model) []Location {
	result := make([]Location, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(Location)
		}
	}
	return result
}
