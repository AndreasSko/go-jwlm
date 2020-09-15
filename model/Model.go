package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/mitchellh/go-wordwrap"
)

// Model represents a general table of the JW Library database and
// defines methods that are needed so entries are mergeable.
type Model interface {
	ID() int
	SetID(int)
	UniqueKey() string
	Equals(m2 Model) bool
	PrettyPrint(db *Database) string
	tableName() string
	idName() string
	scanRow(row *sql.Rows) (Model, error)
}

// MakeModelSlice converts a slice of pointers of model-implementing structs to []model
func MakeModelSlice(arg interface{}) ([]Model, error) {
	slice := reflect.ValueOf(arg)

	if slice.Kind() != reflect.Kind(reflect.Slice) {
		return nil, fmt.Errorf("Can't create []model out of %T", arg)
	}

	c := slice.Len()
	result := make([]Model, c)
	for i := 0; i < c; i++ {
		result[i] = slice.Index(i).Interface().(Model)
	}
	return result, nil
}

// prettyPrint prints the given fields of a Model as a table. If the field
// is empty, its omitted.
func prettyPrint(m Model, fields []string) string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)

Loop:
	for _, fieldName := range fields {
		field := reflect.ValueOf(m).Elem().FieldByName(fieldName)
		if !field.IsValid() {
			panic(fmt.Sprintf("Given struct does not contain field %s", fieldName))
		}
		switch field.Interface().(type) {
		case string:
			fmt.Fprintf(w, "\n%s:\t%s", fieldName, strings.ReplaceAll(wordwrap.WrapString(field.String(), 70), "\n", "\n\t"))
		case sql.NullString:
			if field.Field(1).Bool() == false {
				continue Loop
			}
			fmt.Fprintf(w, "\n%s:\t%s", fieldName, strings.ReplaceAll(wordwrap.WrapString(field.Field(0).String(), 70), "\n", "\n\t"))
		case int:
			fmt.Fprintf(w, "\n%s:\t%d", fieldName, field.Int())
		case sql.NullInt32:
			if field.Field(1).Bool() == false {
				continue Loop
			}
			fmt.Fprintf(w, "\n%s:\t%d", fieldName, field.Field(0).Int())
		default:
			panic(fmt.Sprintf("Unsupported type for field %s", fieldName))
		}
	}
	w.Flush()
	return buf.String()
}
