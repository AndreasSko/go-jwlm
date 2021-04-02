// +build !js

package model

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	// Register SQLite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type fsPersistence struct{}

func getFsPersistence() Persistence {
	pers := fsPersistence{}
	return &pers
}

func getJsPersistence() Persistence {
	panic("getJsPersistence call in non-js runtime")
}

func (pers *fsPersistence) CreateTempStorage(prefix string) (path string, err error) {
	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", prefix)
	if err != nil {
		return "", errors.Wrap(err, "Error while creating temporary directory")
	}
	return tmp, nil
}

func (pers *fsPersistence) StoreSQLiteDB(filename string, dbData []byte) (fullFileName string, err error) {
	return "", errors.Errorf("Not needed to store the JWLBackup when using fsPersistence")
}

func (pers *fsPersistence) OpenSQLiteDB(fullFileName string) (*sql.DB, error) {
	return sql.Open("sqlite3", fullFileName)
}

func (pers *fsPersistence) RetrieveSQLiteData(fullFileName string) ([]byte, error) {
	_, data, err := pers.GetFile(fullFileName)
	return data, err
}

func (pers *fsPersistence) StoreJWLBackup(fullFileName string, archiveData []byte) error {
	return errors.Errorf("Not needed to store the JWLBackup when using fsPersistence")
}

func (pers *fsPersistence) ProcessJWLBackup(fullFileName string, exportPath string) error {

	r, err := zip.OpenReader(fullFileName)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fileReader, err := file.Open()
		if err != nil {
			return errors.Wrap(err, "Error while opening zip file")
		}
		defer fileReader.Close()

		path := filepath.Join(exportPath, file.Name)
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return errors.Wrap(err, "Error while uncompressing zip file")
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return errors.Wrap(err, "Error while copying files from backup to folder")
		}
	}

	return nil
}

func (pers *fsPersistence) GetFile(fullFileName string) (filename string, data []byte, err error) {
	file, err := os.Open(fullFileName)
	if err != nil {
		return "", nil, errors.Wrap(err, fmt.Sprintf("Error opening file at %v", fullFileName))
	}
	defer file.Close()

	blob, err := ioutil.ReadAll(file)
	if err != nil {
		return "", nil, errors.Wrap(err, fmt.Sprintf("Error reading file at %v", fullFileName))
	}

	return file.Name(), blob, nil

}

func (pers *fsPersistence) WriteFile(fullFileName string, data []byte) error {
	if err := ioutil.WriteFile(fullFileName, data, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while saving file at %v", fullFileName))
	}
	return nil
}

func (pers *fsPersistence) CleanupPath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error clearing path %v", path))
	}

	return nil
}
