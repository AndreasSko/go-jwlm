package model

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_getSliceCapacity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := mock.NewRows([]string{"TagId"}).AddRow(3)
	mock.ExpectQuery("SELECT TagId FROM Tag ORDER BY TagId DESC LIMIT 1").WillReturnRows(rows)

	res, err := getSliceCapacity(db, &Tag{})
	assert.NoError(t, err)
	assert.Equal(t, 4, res)

	// Test with empty DB
	rows = mock.NewRows([]string{"TagId"})
	mock.ExpectQuery("SELECT TagId FROM Tag ORDER BY TagId DESC LIMIT 1").WillReturnRows(rows)
	res, err = getSliceCapacity(db, &Tag{})
	assert.NoError(t, err)
	assert.Equal(t, 1, res)
}

func Test_fetchFromSQLite(t *testing.T) {
	path := filepath.Join("testdata", "user_data.db")
	sqlite, err := sql.Open("sqlite3", path+"?immutable=1")
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error while opening SQLite database"))
	}
	defer sqlite.Close()

	blockRange, err := fetchFromSQLite(sqlite, &BlockRange{})
	assert.NoError(t, err)
	assert.Len(t, blockRange, 5)
	assert.Equal(t, &BlockRange{3, 2, 13, sql.NullInt32{Int32: 0, Valid: true}, sql.NullInt32{Int32: 14, Valid: true}, 3}, blockRange[3])

	bookmark, err := fetchFromSQLite(sqlite, &Bookmark{})
	assert.NoError(t, err)
	assert.Len(t, bookmark, 3)
	assert.Equal(t, &Bookmark{2, 3, 7, 4, "Philippians 4", sql.NullString{String: "12 I know how to be low on provisions and how to have an abundance. In everything and in all circumstances I have learned the secret of both how to be full and how to hunger, both how to have an abundance and how to do without. ", Valid: true}, 0, sql.NullInt32{}}, bookmark[2])

	location, err := fetchFromSQLite(sqlite, &Location{})
	assert.NoError(t, err)
	assert.Len(t, location, 8)
	assert.Equal(t, &Location{4, sql.NullInt32{Int32: 66, Valid: true}, sql.NullInt32{Int32: 21, Valid: true}, sql.NullInt32{}, sql.NullInt32{}, 0, sql.NullString{String: "nwtsty", Valid: true}, 2, 0, sql.NullString{String: "Offenbarung 21", Valid: true}}, location[4])

	note, err := fetchFromSQLite(sqlite, &Note{})
	assert.NoError(t, err)
	assert.Len(t, note, 3)
	assert.Equal(t, &Note{2, "F75A18EE-FC17-4E0B-ABB6-CC16DABE9610", sql.NullInt32{Int32: 3, Valid: true}, sql.NullInt32{Int32: 3, Valid: true}, sql.NullString{String: "For all things I have the strength through the one who gives me power.", Valid: true}, sql.NullString{String: "", Valid: true}, "2020-04-14T18:42:14+00:00", 2, sql.NullInt32{Int32: 13, Valid: true}}, note[2])

	tag, err := fetchFromSQLite(sqlite, &Tag{})
	assert.NoError(t, err)
	assert.Len(t, tag, 3)
	assert.Equal(t, &Tag{2, 1, "Strengthening", sql.NullString{}}, tag[2])

	tagMap, err := fetchFromSQLite(sqlite, &TagMap{})
	assert.NoError(t, err)
	assert.Len(t, tagMap, 3)
	assert.Equal(t, &TagMap{2, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 2, Valid: true}, 2, 1}, tagMap[2])

	userMark, err := fetchFromSQLite(sqlite, &UserMark{})
	assert.NoError(t, err)
	assert.Len(t, userMark, 5)
	assert.Equal(t, &UserMark{2, 1, 2, 0, "2C5E7B4A-4997-4EDA-9CFF-38A7599C487B", 1}, userMark[2])
}

func TestDatabase_importSQLite(t *testing.T) {
	db := Database{}

	path := filepath.Join("testdata", "user_data.db")
	assert.NoError(t, db.importSQLite(path))

	// As we already test the correctness in Test_fetchFromSQLite,
	// it should be sufficient to just double-check the size of the slices.
	assert.Len(t, db.BlockRange, 5)
	assert.Len(t, db.Bookmark, 3)
	assert.Len(t, db.Location, 8)
	assert.Len(t, db.Note, 3)
	assert.Len(t, db.Tag, 3)
	assert.Len(t, db.TagMap, 3)
	assert.Len(t, db.UserMark, 5)
}

func TestDatabase_ImportJWLBackup(t *testing.T) {
	db := Database{}

	path := filepath.Join("testdata", "backup.jwlibrary")
	assert.NoError(t, db.ImportJWLBackup(path))

	// As we already test the correctness in Test_fetchFromSQLite,
	// it should be sufficient to just double-check the size of the slices.
	assert.Len(t, db.BlockRange, 5)
	assert.Len(t, db.Bookmark, 3)
	assert.Len(t, db.Location, 8)
	assert.Len(t, db.Note, 3)
	assert.Len(t, db.Tag, 3)
	assert.Len(t, db.TagMap, 3)
	assert.Len(t, db.UserMark, 5)
}

func TestDatabase_ExportJWLBackup(t *testing.T) {
	// Create tmp folder and place all files there
	testFolder := ".jwlm-tmp_test"
	err := os.Mkdir(testFolder, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(testFolder)

	// Test if import->export->import tweakes Data in wrong way
	db := Database{}

	path := filepath.Join("testdata", "backup.jwlibrary")
	assert.NoError(t, db.ImportJWLBackup(path))

	newPath := filepath.Join(testFolder, "backup.jwlibrary")
	assert.NoError(t, db.ExportJWLBackup(newPath))

	db = Database{}
	assert.NoError(t, db.ImportJWLBackup(newPath))

	assert.Len(t, db.BlockRange, 5)
	assert.Equal(t, &BlockRange{3, 2, 13, sql.NullInt32{Int32: 0, Valid: true}, sql.NullInt32{Int32: 14, Valid: true}, 3}, db.BlockRange[3])

	assert.Len(t, db.Bookmark, 3)
	assert.Equal(t, &Bookmark{2, 3, 7, 4, "Philippians 4", sql.NullString{String: "12 I know how to be low on provisions and how to have an abundance. In everything and in all circumstances I have learned the secret of both how to be full and how to hunger, both how to have an abundance and how to do without. ", Valid: true}, 0, sql.NullInt32{}}, db.Bookmark[2])

	assert.Len(t, db.Location, 8)
	assert.Equal(t, &Location{4, sql.NullInt32{Int32: 66, Valid: true}, sql.NullInt32{Int32: 21, Valid: true}, sql.NullInt32{}, sql.NullInt32{}, 0, sql.NullString{String: "nwtsty", Valid: true}, 2, 0, sql.NullString{String: "Offenbarung 21", Valid: true}}, db.Location[4])

	assert.Len(t, db.Note, 3)
	assert.Equal(t, &Note{2, "F75A18EE-FC17-4E0B-ABB6-CC16DABE9610", sql.NullInt32{Int32: 3, Valid: true}, sql.NullInt32{Int32: 3, Valid: true}, sql.NullString{String: "For all things I have the strength through the one who gives me power.", Valid: true}, sql.NullString{String: "", Valid: true}, "2020-04-14T18:42:14+00:00", 2, sql.NullInt32{Int32: 13, Valid: true}}, db.Note[2])

	assert.Len(t, db.Tag, 3)
	assert.Equal(t, &Tag{2, 1, "Strengthening", sql.NullString{}}, db.Tag[2])

	assert.Len(t, db.TagMap, 3)
	assert.Equal(t, &TagMap{2, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 2, Valid: true}, 2, 1}, db.TagMap[2])

	assert.Len(t, db.UserMark, 5)
	assert.Equal(t, &UserMark{2, 1, 2, 0, "2C5E7B4A-4997-4EDA-9CFF-38A7599C487B", 1}, db.UserMark[2])
}

func Test_createEmptySQLiteDB(t *testing.T) {
	// Create tmp folder and place all files there
	err := os.Mkdir(tmpFolder, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(tmpFolder)

	path := filepath.Join(tmpFolder, "user_data.db")
	err = createEmptySQLiteDB(path)
	assert.NoError(t, err)

	// Test if file has correct hash
	f, err := os.Open(path)
	if err != nil {
		errors.Wrap(err, "Error while opening SQLite file to calculate hash")
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	assert.Equal(t, "150423e70425df0b5c2fe76d445fa0b488f73114cdb647abd4deff52aa4f9159", hash)
}

func TestDatabase_saveToNewSQLite(t *testing.T) {
	// Create tmp folder and place all files there
	err := os.Mkdir(tmpFolder, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(tmpFolder)

	db := Database{
		BlockRange: []*BlockRange{{3, 2, 13, sql.NullInt32{Int32: 0, Valid: true}, sql.NullInt32{Int32: 14, Valid: true}, 3}},
		Bookmark:   []*Bookmark{{2, 3, 7, 4, "Philippians 4", sql.NullString{String: "12 I know how to be low on provisions and how to have an abundance. In everything and in all circumstances I have learned the secret of both how to be full and how to hunger, both how to have an abundance and how to do without. ", Valid: true}, 0, sql.NullInt32{}}},
		Location:   []*Location{{4, sql.NullInt32{Int32: 66, Valid: true}, sql.NullInt32{Int32: 21, Valid: true}, sql.NullInt32{}, sql.NullInt32{}, 0, sql.NullString{String: "nwtsty", Valid: true}, 2, 0, sql.NullString{String: "Offenbarung 21", Valid: true}}},
		Note:       []*Note{{2, "F75A18EE-FC17-4E0B-ABB6-CC16DABE9610", sql.NullInt32{Int32: 3, Valid: true}, sql.NullInt32{Int32: 3, Valid: true}, sql.NullString{String: "For all things I have the strength through the one who gives me power.", Valid: true}, sql.NullString{String: "!", Valid: true}, "2020-04-14T18:42:14+00:00", 2, sql.NullInt32{Int32: 13, Valid: true}}},
		Tag:        []*Tag{{2, 1, "Strengthening", sql.NullString{}}},
		TagMap:     []*TagMap{{2, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 2, Valid: true}, 2, 1}},
		UserMark:   []*UserMark{{2, 1, 2, 0, "2C5E7B4A-4997-4EDA-9CFF-38A7599C487B", 1}},
	}
	path := filepath.Join(tmpFolder, "user_data.db")
	err = db.saveToNewSQLite(path)
	assert.NoError(t, err)

	db2 := Database{}
	err = db2.importSQLite(path)
	assert.NoError(t, err)

	assert.Equal(t, db.BlockRange[0], db2.BlockRange[3])
	assert.Equal(t, db.Bookmark[0], db2.Bookmark[2])
	assert.Equal(t, db.Location[0], db2.Location[4])
	assert.Equal(t, db.Note[0], db2.Note[2])
	assert.Equal(t, db.TagMap[0], db2.TagMap[2])
	assert.Equal(t, db.UserMark[0], db2.UserMark[2])
}
