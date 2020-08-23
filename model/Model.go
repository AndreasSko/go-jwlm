package model

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Model represents a general table of the JW Library database and
// defines methods that are needed so entries are mergeable.
type Model interface {
	ID() int
	SetID(int)
	UniqueKey() string
	Equals(m2 Model) bool
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
