package model

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
)

// Tag represents the Tag table inside the JW Library database
type Tag struct {
	TagID         int
	TagType       int
	Name          string
	ImageFilename sql.NullString
}

// ID returns the ID of the entry
func (m *Tag) ID() int {
	return m.TagID
}

// SetID sets the ID of the entry
func (m *Tag) SetID(id int) {
	m.TagID = id
}

// UniqueKey returns the key that makes this Tag unique,
// so it can be used as a key in a map.
func (m *Tag) UniqueKey() string {
	var sb strings.Builder
	sb.Grow(15)
	sb.WriteString(strconv.FormatInt(int64(m.TagType), 10))
	sb.WriteString("_")
	sb.WriteString(m.Name)
	return sb.String()
}

// Equals checks if the Tag is equal to the given one. The
// check won't include the TagID.
func (m *Tag) Equals(m2 Model) bool {
	if m2, ok := m2.(*Tag); ok {
		return m.TagType == m2.TagType &&
			m.Name == m2.Name &&
			m.ImageFilename == m2.ImageFilename
	}

	return false
}

// RelatedEntries returns entries that are related to this one
func (m *Tag) RelatedEntries(db *Database) []Model {
	return []Model{}
}

// PrettyPrint prints Tag in a human readable format and
// adds information about related entries if helpful.
func (m *Tag) PrettyPrint(db *Database) string {
	fields := []string{"Name"}
	return prettyPrint(m, fields)
}

// MarshalJSON returns the JSON encoding of the entry
func (m Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type          string
		TagID         int
		TagType       int
		Name          string
		ImageFilename sql.NullString
	}{
		Type:          "Tag",
		TagID:         m.TagID,
		TagType:       m.TagType,
		Name:          m.Name,
		ImageFilename: m.ImageFilename,
	})
}

func (m *Tag) tableName() string {
	return "Tag"
}

func (m *Tag) idName() string {
	return "TagId"
}

func (m *Tag) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.TagID, &m.TagType, &m.Name, &m.ImageFilename)
	return m, err
}

// MakeSlice converts a slice of the generice interface model
func (Tag) MakeSlice(mdl []Model) []*Tag {
	result := make([]*Tag, len(mdl))
	for i := range mdl {
		if mdl[i] == nil {
			continue
		}
		tag := mdl[i].(*Tag)

		result[i] = tag
	}
	return result
}
