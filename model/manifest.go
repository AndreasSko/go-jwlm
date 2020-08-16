package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type manifest struct {
	CreationDate   string         `json:"creationDate"`
	UserDataBackup userDataBackup `json:"userDataBackup"`
	Name           string         `json:"name"`
	Type           int            `json:"type"`
	Version        int            `json:"version"`
}
type userDataBackup struct {
	LastModifiedDate string `json:"lastModifiedDate"`
	Hash             string `json:"hash"`
	DatabaseName     string `json:"databaseName"`
	SchemaVersion    int    `json:"schemaVersion"`
	DeviceName       string `json:"deviceName"`
}

// validateManifest checks if the backup file is compatible by validating its manifest.json
func validateManifest(path string) error {
	const version = 1
	const schemaVersion = 8

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	blob, _ := ioutil.ReadAll(file)

	var manifest manifest
	err = json.Unmarshal([]byte(blob), &manifest)
	if err != nil {
		return errors.Wrap(err, "Could not unmarshall backup manifest file")
	}

	if manifest.Version != version {
		return fmt.Errorf("Manifest version is incompatible. Should be %d is %d", version, manifest.Version)
	}

	if manifest.UserDataBackup.SchemaVersion != schemaVersion {
		return fmt.Errorf("Schema version is incompatible. Should be %d is %d", schemaVersion, manifest.UserDataBackup.SchemaVersion)
	}

	return nil
}

// createManifest creates a manifest for a given SQLite database and saves it at path
func createManifest(backupName string, dbFile string, path string) error {
	// Get SHA256 of SQLite file
	f, err := os.Open(dbFile)
	if err != nil {
		errors.Wrap(err, "Error while opening SQLite file to calculate hash")
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	udb := userDataBackup{
		LastModifiedDate: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Hash:             hash,
		DatabaseName:     filepath.Base(dbFile),
		SchemaVersion:    8,
		DeviceName:       "go-jwlm",
	}

	manifest := manifest{
		CreationDate:   time.Now().Format("2006-01-02"),
		UserDataBackup: udb,
		Name:           backupName,
		Type:           0,
		Version:        1,
	}

	bytes, err := json.Marshal(manifest)
	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while saving manifest file to %v", path))
	}

	return nil
}
