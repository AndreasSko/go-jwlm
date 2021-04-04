package storage

import (
	"archive/zip"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
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
	mfstBytes, err := readFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading manifest")
	}
	// Parse manifest
	mfst := model.Manifest{}
	err = json.Unmarshal([]byte(mfstBytes), &mfst) // TODO: Is []byte() needed??
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshall backup manifest file")
	}

	// Make sure that we support this backup version
	if err := mfst.ValidateManifest(); err != nil {
		return nil, err
	}

	// Open SQLite
	dbFilename := filepath.Join(tmp, mfst.UserDataBackup.DatabaseName)
	sqlite, err := sql.Open("sqlite3", dbFilename+"?immutable=1")
	if err != nil {
		return nil, errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	db := &model.Database{}
	err = db.Import(sqlite)
	return db, err
}

func ExportJWLBackup(db *model.Database, filename string) error {
	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", "go-jwlm")
	if err != nil {
		return errors.Wrap(err, "Error while creating temporary directory")
	}
	//defer os.RemoveAll(tmp)

	// Create user_data.db
	dbPath := filepath.Join(tmp, "user_data.db")
	if err := createEmptySQLiteDB(dbPath); err != nil {
		return errors.Wrap(err, "Error while creating new empty SQLite database")
	}

	sqlite, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	db.Export(sqlite)

	mfstHash, err := hashOfFile(dbPath)
	if err != nil {
		return err
	}

	mfst := model.GenerateManifest("go-jwlm", "user_data.db", mfstHash)
	bytes, err := json.Marshal(mfst)
	if err != nil {
		return errors.Wrap(err, "Error while marshalling manifest")
	}
	mfstPath := filepath.Join(tmp, "manifest.json")
	if err := ioutil.WriteFile(mfstPath, bytes, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while saving manifest file at %v", mfstPath))
	}

	// Store files in .jwlibrary (zip)-file
	files := []string{dbPath, mfstPath}
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

func hashOfFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "Error while opening file %s to calculate hash", path)
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
