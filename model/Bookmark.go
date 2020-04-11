package model

import "database/sql"

// Bookmark represents the Bookmark table inside the JW Library database
type Bookmark struct {
	BookmarkID            int
	LocationID            int
	PublicationLocationID int
	Slot                  int
	Title                 string
	Snippet               sql.NullString
	BlockType             int
	BlockIdentifier       sql.NullInt32
}

// ID returns the ID of the entry
func (m Bookmark) ID() int {
	return m.BookmarkID
}

func (m Bookmark) tableName() string {
	return "Bookmark"
}

func (m Bookmark) idName() string {
	return "BookmarkId"
}

func (m Bookmark) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.BookmarkID, &m.LocationID, &m.PublicationLocationID, &m.Slot, &m.Title,
		&m.Snippet, &m.BlockType, &m.BlockIdentifier)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (Bookmark) makeSlice(mdl []model) []Bookmark {
	result := make([]Bookmark, len(mdl))
	for i := range mdl {
		if mdl[i] != nil {
			result[i] = mdl[i].(Bookmark)
		}
	}
	return result
}
