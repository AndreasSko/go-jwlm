package model

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	// Register SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// Destination to temporarily unzip backups
const tmpFolder = ".jwlm-tmp"
const dbFilename = "user_data.db"
const manifestFilename = "manifest.json"

// Database represents the JW Library database as a struct
type Database struct {
	BlockRange []BlockRange
	Bookmark   []Bookmark
	Location   []Location
	Note       []Note
	Tag        []Tag
	TagMap     []TagMap
	UserMark   []UserMark
}

// ImportJWLBackup unzips a given JW Library Backup file and imports the
// included SQLite DB to the Database struct
func (db *Database) ImportJWLBackup(filename string) error {
	// Create tmp folder and unzip backup content there
	if err := os.Mkdir(tmpFolder, 0755); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not create temporary directory %s", tmpFolder))
	}
	defer os.RemoveAll(tmpFolder)

	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		path := filepath.Join(tmpFolder, file.Name)
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return errors.Wrap(err, "Error while copying files from backup to temporary folder")
		}
	}

	// Make sure that we support this backup version
	path := filepath.Join(tmpFolder, manifestFilename)
	if err := validateManifest(path); err != nil {
		return err
	}

	// Fill the Database with actual data
	path = filepath.Join(tmpFolder, dbFilename)
	return db.importSQLite(path)
}

// importSQLite imports a given SQLite DB into the Database struct
func (db *Database) importSQLite(filename string) error {
	// Open SQLite file as immutable to avoid locks (and therefore speed up import)
	sqlite, err := sql.Open("sqlite3", filename+"?immutable=1")
	if err != nil {
		return errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	// Fill each table struct separately (did not find a DRYer solution yet..)
	mdl, err := fetchFromSQLite(sqlite, &BlockRange{})
	if err != nil {
		return err
	}
	db.BlockRange = BlockRange{}.makeSlice(mdl)

	mdl, err = fetchFromSQLite(sqlite, &Bookmark{})
	if err != nil {
		return err
	}
	db.Bookmark = Bookmark{}.makeSlice(mdl)

	mdl, err = fetchFromSQLite(sqlite, &Location{})
	if err != nil {
		return err
	}
	db.Location = Location{}.makeSlice(mdl)

	mdl, err = fetchFromSQLite(sqlite, &Note{})
	if err != nil {
		return err
	}
	db.Note = Note{}.makeSlice(mdl)

	mdl, err = fetchFromSQLite(sqlite, &Tag{})
	if err != nil {
		return err
	}
	db.Tag = Tag{}.makeSlice(mdl)

	mdl, err = fetchFromSQLite(sqlite, &TagMap{})
	if err != nil {
		return err
	}
	db.TagMap = TagMap{}.makeSlice(mdl)

	mdl, err = fetchFromSQLite(sqlite, &UserMark{})
	if err != nil {
		return err
	}

	db.UserMark = UserMark{}.makeSlice(mdl)
	if err != nil {
		return err
	}

	return nil
}

// fetchFromSQLite fetches the entries for a given modelType and returns a slice
// of entries, for which the index corresponds to the ID in the SQLite DB
func fetchFromSQLite(sqlite *sql.DB, modelType model) ([]model, error) {
	// Create slice of correct size (number of entries)
	capacity, err := getSliceCapacity(sqlite, modelType)
	if err != nil {
		return nil, errors.Wrap(err, "Could not determine number of entries in SQLite database")
	}
	result := make([]model, capacity)

	rows, err := sqlite.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY %s", modelType.tableName(), modelType.idName()))
	if err != nil {
		return nil, errors.Wrap(err, "Error while querying SQLite database")
	}

	// Put entries in slice with the index coresponding to the ID in the SQLite DB
	defer rows.Close()
	for rows.Next() {
		var m model
		switch tp := modelType.(type) {
		case *BlockRange:
			m = BlockRange{}
		case *Bookmark:
			m = Bookmark{}
		case *Location:
			m = Location{}
		case *Note:
			m = Note{}
		case *Tag:
			m = Tag{}
		case *TagMap:
			m = TagMap{}
		case *UserMark:
			m = UserMark{}
		default:
			panic(fmt.Sprintf("Fetching %T is not supported!", tp))
		}
		mn, err := m.scanRow(rows)
		if err != nil {
			return nil, errors.Wrap(err, "Error while scanning results from SQLite database")
		}
		result[mn.ID()] = mn
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error while scanning results from SQLite database")
	}

	return result, nil
}

// getSliceCapacity determines the needed capacity for a slice from a table
// by looking at the highest ID in the DB
func getSliceCapacity(sqlite *sql.DB, modelType model) (int, error) {
	row, err := sqlite.Query(fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC LIMIT 1",
		modelType.idName(), modelType.tableName(), modelType.idName()))
	if err != nil {
		return 0, err
	}
	defer row.Close()

	capacity := 0
	for row.Next() {
		if err := row.Scan(&capacity); err != nil {
			return 0, err
		}
	}

	// Index in DB starts with 1, so 0 is always nil
	return capacity + 1, nil
}
