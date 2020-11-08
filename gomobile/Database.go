package gomobile

import (
	"errors"

	"github.com/AndreasSko/go-jwlm/model"
)

// DatabaseWrapper wraps the left, right, and merged
// Database structs so they can be used with Gomobile.
type DatabaseWrapper struct {
	left   *model.Database
	right  *model.Database
	merged *model.Database

	// Temporary databases allow to run a merge function multiple
	// times without changing the content of the original databases.
	leftTmp  *model.Database
	rightTmp *model.Database
}

// ImportJWLBackup imports a .jwlibrary backup file into the struct
// on the given side.
func (dbw *DatabaseWrapper) ImportJWLBackup(filename string, side string) error {
	db := &model.Database{}

	if err := db.ImportJWLBackup(filename); err != nil {
		return err
	}

	switch side {
	case "leftSide":
		dbw.left = db
	case "rightSide":
		dbw.right = db
	default:
		return errors.New("Only leftSide and rightSide are valid for importing backups")
	}

	return nil
}

// Init initializes the DatabaseWrapper to prepare for subsequent
// function calls. Should be called after ImportJWLBackup.
func (dbw *DatabaseWrapper) Init() {
	dbw.leftTmp = model.MakeDatabaseCopy(dbw.left)
	dbw.rightTmp = model.MakeDatabaseCopy(dbw.right)
	dbw.merged = &model.Database{}
}

// DBIsLoaded indicates if a DB on the given side has been loaded.
func (dbw *DatabaseWrapper) DBIsLoaded(side string) bool {
	switch side {
	case "leftSide":
		return dbw.left != nil
	case "rightSide":
		return dbw.right != nil
	case "mergeSide":
		return dbw.merged != nil
	}

	return false
}

// ExportMerged exports the merged database to filename.
func (dbw *DatabaseWrapper) ExportMerged(filename string) error {
	return dbw.merged.ExportJWLBackup(filename)
}
