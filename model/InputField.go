package model

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
)

// InputField represents the InputField table inside the JW Library database
type InputField struct {
	LocationID int
	TextTag    string
	Value      string
	pseudoID   int
}

// ID returns the ID of the entry. As the InputField table does not have
// an ID, we are using a pseudoID, so the rest of the merge logic is
// still able to run as usual.
func (m *InputField) ID() int {
	return m.pseudoID
}

// SetID sets the ID of the entry. As the InputField table does not have
// an ID, this function does nothing.
func (m *InputField) SetID(id int) {
}

// UniqueKey returns the key that makes this InputField unique,
// so it can be used as a key in a map.
func (m *InputField) UniqueKey() string {
	var sb strings.Builder
	sb.Grow(15)
	sb.WriteString(strconv.FormatInt(int64(m.LocationID), 10))
	sb.WriteString("_")
	sb.WriteString(m.TextTag)
	return sb.String()
}

// Equals checks if the InputField is equal to the given one.
func (m *InputField) Equals(m2 Model) bool {
	if m2, ok := m2.(*InputField); ok {
		return m.LocationID == m2.LocationID &&
			m.TextTag == m2.TextTag &&
			m.Value == m2.Value
	}

	return false
}

// RelatedEntries returns entries that are related to this one.
func (m *InputField) RelatedEntries(db *Database) Related {
	result := Related{}

	if location := db.FetchFromTable("Location", int(m.LocationID)); location != nil {
		result.Location = location.(*Location)
	}

	return result
}

// PrettyPrint prints InputField in a human readable format and
// adds information about related entries if helpful.
func (m *InputField) PrettyPrint(db *Database) string {
	var result string

	if location := db.FetchFromTable("Location", m.LocationID); location != nil {
		result += location.PrettyPrint(db) + "\n"
	}

	result += prettyPrint(m, []string{"TextTag", "Value"})

	return result
}

// MarshalJSON returns the JSON encoding of the entry.
func (m InputField) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type       string `json:"type"`
		LocationID int    `json:"locationId"`
		TextTag    string `json:"textTag"`
		Value      string `json:"value"`
	}{
		Type:       "InputField",
		LocationID: m.LocationID,
		TextTag:    m.TextTag,
		Value:      m.Value,
	})
}

func (m *InputField) tableName() string {
	return "InputField"
}

func (m *InputField) idName() string {
	return ""
}

func (m *InputField) scanRow(rows *sql.Rows) (Model, error) {
	err := rows.Scan(&m.LocationID, &m.TextTag, &m.Value)
	return m, err
}

// MakeSlice converts a slice of the generice interface model.
func (InputField) MakeSlice(mdl []Model) []*InputField {
	result := make([]*InputField, len(mdl))
	for i := range mdl {
		if mdl[i] == nil {
			continue
		}
		inputField := mdl[i].(*InputField)

		result[i] = inputField
	}
	return result
}
