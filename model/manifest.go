package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type manifest struct {
	CreationDate   string `json:"creationDate"`
	UserDataBackup struct {
		LastModifiedDate string `json:"lastModifiedDate"`
		Hash             string `json:"hash"`
		DatabaseName     string `json:"databaseName"`
		SchemaVersion    int    `json:"schemaVersion"`
		DeviceName       string `json:"deviceName"`
	} `json:"userDataBackup"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Version int    `json:"version"`
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
