package gomobile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/tj/assert"
)

var backupFile = filepath.Join("..", "model", "testdata", "backup.jwlibrary")

func TestDatabaseWrapper_ImportJWLBackup(t *testing.T) {
	dbw := &DatabaseWrapper{}

	assert.Error(t, dbw.ImportJWLBackup("wrongFile", "leftSide"))
	assert.NoError(t, dbw.ImportJWLBackup(backupFile, "leftSide"))
	assert.NoError(t, dbw.ImportJWLBackup(backupFile, "rightSide"))
	assert.EqualError(t, dbw.ImportJWLBackup(backupFile, "wrongSide"), "Only leftSide and rightSide are valid for importing backups")

	assert.Len(t, dbw.left.BlockRange, 5)
	assert.Len(t, dbw.left.Bookmark, 3)
	assert.Len(t, dbw.left.Location, 8)
	assert.Len(t, dbw.left.Note, 3)
	assert.Len(t, dbw.left.Tag, 3)
	assert.Len(t, dbw.left.TagMap, 3)
	assert.Len(t, dbw.left.UserMark, 5)

	assert.True(t, dbw.left.Equals(dbw.right))
}

func TestDatabaseWrapper_Init(t *testing.T) {
	dbw := &DatabaseWrapper{}

	assert.NoError(t, dbw.ImportJWLBackup(backupFile, "leftSide"))
	assert.NoError(t, dbw.ImportJWLBackup(backupFile, "rightSide"))

	dbw.Init()

	assert.True(t, dbw.merged.Equals(&model.Database{}))
	assert.True(t, dbw.leftTmp.Equals(dbw.left))
	assert.True(t, dbw.rightTmp.Equals(dbw.right))

	dbw.leftTmp.Bookmark[1].Title = "Tweaked"
	assert.False(t, dbw.leftTmp.Equals(dbw.left))

	dbw.merged = model.MakeDatabaseCopy(dbw.left)
	dbw.Init()
	assert.True(t, dbw.merged.Equals(&model.Database{}))
	assert.True(t, dbw.leftTmp.Equals(dbw.left))
	assert.True(t, dbw.rightTmp.Equals(dbw.right))
}

func TestDatabaseWrapper_DBIsLoaded(t *testing.T) {
	dbw := &DatabaseWrapper{}

	assert.False(t, dbw.DBIsLoaded("leftSide"))
	assert.False(t, dbw.DBIsLoaded("rightSide"))
	assert.False(t, dbw.DBIsLoaded("mergeSide"))

	assert.NoError(t, dbw.ImportJWLBackup(backupFile, "leftSide"))
	assert.True(t, dbw.DBIsLoaded("leftSide"))
	assert.False(t, dbw.DBIsLoaded("rightSide"))
	assert.False(t, dbw.DBIsLoaded("mergeSide"))

	assert.NoError(t, dbw.ImportJWLBackup(backupFile, "rightSide"))
	assert.True(t, dbw.DBIsLoaded("leftSide"))
	assert.True(t, dbw.DBIsLoaded("rightSide"))
	assert.False(t, dbw.DBIsLoaded("mergeSide"))

	dbw.merged = model.MakeDatabaseCopy(dbw.left)
	assert.True(t, dbw.DBIsLoaded("mergeSide"))

	assert.False(t, dbw.DBIsLoaded("wrongSide"))
}

func TestDatabaseWrapper_ExportMerged(t *testing.T) {
	tmp, err := ioutil.TempDir("", "go-jwlm")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)

	dbw := &DatabaseWrapper{}
	dbw.merged = &model.Database{}

	assert.NoError(t, dbw.merged.ImportJWLBackup(backupFile))

	newBackup := filepath.Join(tmp, "test.jwlibrary")
	assert.NoError(t, dbw.ExportMerged(newBackup))

	newDB := &model.Database{}
	assert.NoError(t, newDB.ImportJWLBackup(newBackup))
	assert.True(t, dbw.merged.Equals(newDB))
}
