package model

import "database/sql"

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

func (m *Tag) tableName() string {
	return "Tag"
}

func (m *Tag) idName() string {
	return "TagId"
}

func (m *Tag) scanRow(rows *sql.Rows) (model, error) {
	err := rows.Scan(&m.TagID, &m.TagType, &m.Name, &m.ImageFilename)
	return m, err
}

// makeSlice converts a slice of the generice interface model
func (Tag) makeSlice(mdl []*model) []*Tag {
	result := make([]*Tag, len(mdl))
	for i := range mdl {
		if mdl[i] == nil {
			continue
		}
		tag := (*mdl[i]).(*Tag)

		// The "Favorite" tag is already included with a fresh JW-Library installation
		if tag.TagID == 1 && tag.TagType == 0 && tag.Name == "Favorite" {
			continue
		}

		result[i] = tag
	}
	return result
}
