package model

import (
	"path/filepath"
	"testing"
)

func Test_validateManifest(t *testing.T) {
	path := filepath.Join("testdata", "manifest_correct.json")
	if err := validateManifest(path); err != nil {
		t.Error("Manifest validation should have been successful")
	}

	path = filepath.Join("testdata", "manifest_oudated.json")
	if err := validateManifest(path); err == nil {
		t.Error("Manifest validation should have failed")
	}
}
