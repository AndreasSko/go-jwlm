package model

import (
	"database/sql"
	"fmt"
	"reflect"
)

type model interface {
	ID() int
	tableName() string
	idName() string
	scanRow(row *sql.Rows) (model, error)
}

// makeModelSclie converts a slice of model-implementing structs to []model
func makeModelSlice(arg interface{}) ([]model, error) {
	slice := reflect.ValueOf(arg)

	if slice.Kind() != reflect.Kind(reflect.Slice) {
		return nil, fmt.Errorf("Can't create []model out of %T", arg)
	}

	c := slice.Len()
	result := make([]model, c)
	for i := 0; i < c; i++ {
		result[i] = slice.Index(i).Interface().(model)
	}
	return result, nil
}
