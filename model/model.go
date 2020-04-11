package model

import "database/sql"

type model interface {
	ID() int
	tableName() string
	idName() string
	scanRow(row *sql.Rows) (model, error)
}
