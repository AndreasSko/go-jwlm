package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
)

type Manifest struct {
	CreationDate   string         `json:"creationDate"`
	UserDataBackup UserDataBackup `json:"userDataBackup"`
	Name           string         `json:"name"`
	Type           int            `json:"type"`
	Version        int            `json:"version"`
}
type UserDataBackup struct {
	LastModifiedDate string `json:"lastModifiedDate"`
	Hash             string `json:"hash"`
	DatabaseName     string `json:"databaseName"`
	SchemaVersion    int    `json:"schemaVersion"`
	DeviceName       string `json:"deviceName"`
}

// importManifest imports a manifest.json at path
func (mfst *Manifest) importManifest(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	blob, _ := ioutil.ReadAll(file)

	err = json.Unmarshal([]byte(blob), &mfst)
	if err != nil {
		return errors.Wrap(err, "Could not unmarshall backup manifest file")
	}

	return nil
}

// ValidateManifest checks if the backup file is compatible by validating the manifest
func (mfst *Manifest) ValidateManifest() error {
	const version = 1
	const schemaVersion = 8

	if mfst.Version != version {
		return fmt.Errorf("Manifest version is incompatible. Should be %d is %d. "+
			"You might need to upgrade to a newer version of JW Library first", version, mfst.Version)
	}

	if mfst.UserDataBackup.SchemaVersion != schemaVersion {
		return fmt.Errorf("Schema version is incompatible. Should be %d is %d. "+
			"You might need to upgrade to a newer version of JW Library first", schemaVersion, mfst.UserDataBackup.SchemaVersion)
	}

	return nil
}

// GenerateManifest generates a manifest from the given information, which can
// later be exported
func GenerateManifest(backupName string, dbFilename string, dbHash string) *Manifest {
	return &Manifest{
		CreationDate: time.Now().Format("2006-01-02"),
		UserDataBackup: UserDataBackup{
			LastModifiedDate: time.Now().Format("2006-01-02T15:04:05-07:00"),
			Hash:             dbHash,
			DatabaseName:     dbFilename,
			SchemaVersion:    8,
			DeviceName:       "go-jwlm",
		},
		Name:    backupName,
		Type:    0,
		Version: 1,
	}
}

// exportManifest exports a manifest at path
func (mfst *Manifest) exportManifest(path string) error {
	bytes, err := json.Marshal(mfst)
	if err != nil {
		return errors.Wrap(err, "Error while marshalling manifest")
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while saving manifest file at %v", path))
	}

	return nil
}
