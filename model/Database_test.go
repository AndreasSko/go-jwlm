package model

import (
	"database/sql"
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

	res, err := getSliceCapacity(db, Tag{})
	assert.NoError(t, err)
	assert.Equal(t, 4, res)

	// Test with empty DB
	rows = mock.NewRows([]string{"TagId"})
	mock.ExpectQuery("SELECT TagId FROM Tag ORDER BY TagId DESC LIMIT 1").WillReturnRows(rows)
	res, err = getSliceCapacity(db, Tag{})
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
	assert.Contains(t, blockRange, BlockRange{3, 2, 13, sql.NullInt32{Int32: 0, Valid: true}, sql.NullInt32{Int32: 14, Valid: true}, 3})

	bookmark, err := fetchFromSQLite(sqlite, &Bookmark{})
	assert.NoError(t, err)
	assert.Len(t, bookmark, 3)
	assert.Contains(t, bookmark, Bookmark{2, 3, 7, 4, "Philippians 4", sql.NullString{String: "12 I know how to be low on provisions and how to have an abundance. In everything and in all circumstances I have learned the secret of both how to be full and how to hunger, both how to have an abundance and how to do without. ", Valid: true}, 0, sql.NullInt32{}})

	location, err := fetchFromSQLite(sqlite, &Location{})
	assert.NoError(t, err)
	assert.Len(t, location, 8)
	assert.Contains(t, location, Location{4, sql.NullInt32{Int32: 66, Valid: true}, sql.NullInt32{Int32: 21, Valid: true}, sql.NullInt32{}, sql.NullInt32{}, 0, sql.NullString{String: "nwtsty", Valid: true}, 2, 0, sql.NullString{String: "Offenbarung 21", Valid: true}})

	note, err := fetchFromSQLite(sqlite, &Note{})
	assert.NoError(t, err)
	assert.Len(t, note, 3)
	assert.Contains(t, note, Note{2, "F75A18EE-FC17-4E0B-ABB6-CC16DABE9610", sql.NullInt32{Int32: 3, Valid: true}, sql.NullInt32{Int32: 3, Valid: true}, sql.NullString{String: "For all things I have the strength through the one who gives me power.", Valid: true}, sql.NullString{String: "!", Valid: true}, "2020-04-14T18:42:14+00:00", 2, sql.NullInt32{Int32: 13, Valid: true}})

	tag, err := fetchFromSQLite(sqlite, &Tag{})
	assert.NoError(t, err)
	assert.Len(t, tag, 3)
	assert.Contains(t, tag, Tag{2, 1, "Strengthening", sql.NullString{}})

	tagMap, err := fetchFromSQLite(sqlite, &TagMap{})
	assert.NoError(t, err)
	assert.Len(t, tagMap, 3)
	assert.Contains(t, tagMap, TagMap{2, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 0, Valid: false}, sql.NullInt32{Int32: 2, Valid: true}, 2, 1})

	userMark, err := fetchFromSQLite(sqlite, &UserMark{})
	assert.NoError(t, err)
	assert.Len(t, userMark, 5)
	assert.Contains(t, userMark, UserMark{2, 1, 2, 0, "2C5E7B4A-4997-4EDA-9CFF-38A7599C487B", 1})
}

func TestDatabase_importSQLite(t *testing.T) {
	db := Database{}

	path := filepath.Join("testdata", "user_data.db")
	assert.NoError(t, db.importSQLite(path))

	// As we already test the correct in Test_fetchFromSQLite,
	// it should be sufficient to just double-check the size of the slices
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

	// As we already test the correct in Test_fetchFromSQLite,
	// it should be sufficient to just double-check the size of the slices
	assert.Len(t, db.BlockRange, 5)
	assert.Len(t, db.Bookmark, 3)
	assert.Len(t, db.Location, 8)
	assert.Len(t, db.Note, 3)
	assert.Len(t, db.Tag, 3)
	assert.Len(t, db.TagMap, 3)
	assert.Len(t, db.UserMark, 5)
}
