package storage

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/pkg/errors"

	// Register SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

const manifestFilename = "manifest.json"

func ImportJWLBackup(filename string) (*model.Database, error) {
	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", "go-jwlm")
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating temporary directory")
	}
	defer os.RemoveAll(tmp)

	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, file := range r.File {
		fileReader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fileReader.Close()

		path := filepath.Join(tmp, file.Name)
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return nil, err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return nil, errors.Wrap(err, "Error while copying files from backup to temporary folder")
		}
	}

	// Read manifest
	path := filepath.Join(tmp, manifestFilename)
	mfstByte, err := readFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading manifest")
	}

	// Open SQLite
	sqlite, err := sql.Open("sqlite3", filename+"?immutable=1")
	if err != nil {
		return nil, errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	db := &model.Database{}
	err = db.Import(mfstByte, sqlite)
	return db, err
}

func ExportJWLBackup(db *model.Database, filename string) error {
	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", "go-jwlm")
	if err != nil {
		return errors.Wrap(err, "Error while creating temporary directory")
	}
	defer os.RemoveAll(tmp)

	// Create user_data.db
	dbPath := filepath.Join(tmp, "user_data.db")
	if err := createEmptySQLiteDB(filename); err != nil {
		return errors.Wrap(err, "Error while creating new empty SQLite database")
	}

	sqlite, err := sql.Open("sqlite3", filename)
	if err != nil {
		return errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	db.Export(sqlite)

	// Store files in .jwlibrary (zip)-file
	files := []string{dbPath, manifestPath}
	if err := zipFiles(filename, files); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while storing files in zip archive %s", filename))
	}

	return nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	blob, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

// createEmptySQLiteDB creates a new SQLite database at filename with the base user_data.db from JWLibrary
func createEmptySQLiteDB(filename string) error {
	userData, err := Asset("user_data.db")
	if err != nil {
		return errors.Wrap(err, "Error while fetching user_data.db")
	}

	if err := ioutil.WriteFile(filename, userData, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while saving new SQLite database at %s", filename))
	}

	return nil
}
