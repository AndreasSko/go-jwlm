package model

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"

	// Register SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

const manifestFilename = "manifest.json"

// Database represents the JW Library database as a struct
type Database struct {
	BlockRange []*BlockRange
	Bookmark   []*Bookmark
	InputField []*InputField
	Location   []*Location
	Note       []*Note
	Tag        []*Tag
	TagMap     []*TagMap
	UserMark   []*UserMark
}

// FetchFromTable tries to fetch a entry with the given ID. If it can't find it
// or the entry is empty it returns nil.
func (db *Database) FetchFromTable(tableName string, id int) Model {
	if db == nil {
		return nil
	}

	table := reflect.ValueOf(db).Elem().FieldByName(tableName)
	if !table.IsValid() {
		panic(fmt.Sprintf("Table %s does not exist in Database", tableName))
	}

	if id >= table.Len() {
		return nil
	}
	if table.Index(id).IsNil() {
		return nil
	}

	return table.Index(id).Interface().(Model)
}

// MakeDatabaseCopy creates a deep copy of the given Database, so elements of
// the copy can be safely updated without affecting the original one.
func MakeDatabaseCopy(db *Database) *Database {
	newDB := &Database{}

	dbFields := reflect.ValueOf(db).Elem()
	for i := 0; i < dbFields.NumField(); i++ {
		field := dbFields.Field(i)
		if !field.CanInterface() {
			continue
		}

		tp := field.Kind()
		switch tp {
		case reflect.Slice:
			cpSlice := reflect.MakeSlice(field.Type(), field.Len(), field.Len())

			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)

				if elem.IsNil() {
					continue
				}

				switch t := elem.Interface().(type) {
				case Model:
					cpSlice.Index(j).Set(reflect.ValueOf(MakeModelCopy(elem.Interface().(Model))))
				default:
					panic(fmt.Sprintf("Element type %T is not supported for copying", t))
				}
			}
			cpField := reflect.ValueOf(newDB).Elem().Field(i)
			cpField.Set(cpSlice)
		default:
			panic(fmt.Sprintf("Field type %T is not supported for copying", tp))
		}
	}

	return newDB
}

// Equals checks if all entries of a Database are equal.
func (db *Database) Equals(other *Database) bool {
	// Make copy of DBs so we can safely transform them if necessary
	dbCp := MakeDatabaseCopy(db)
	otherCp := MakeDatabaseCopy(other)

	// Sort all tables by UniqueKey and update IDs in other tables
	for _, db := range []*Database{dbCp, otherCp} {
		locIDChanges := SortByUniqueKey(&db.Location)
		UpdateIDs(db.Bookmark, "LocationID", locIDChanges)
		UpdateIDs(db.Bookmark, "PublicationLocationID", locIDChanges)
		UpdateIDs(db.InputField, "LocationID", locIDChanges)
		UpdateIDs(db.Note, "LocationID", locIDChanges)
		UpdateIDs(db.TagMap, "LocationID", locIDChanges)
		UpdateIDs(db.UserMark, "LocationID", locIDChanges)

		SortByUniqueKey(&db.Bookmark)
		SortByUniqueKey(&db.InputField)

		tagIDChanges := SortByUniqueKey(&db.Tag)
		UpdateIDs(db.TagMap, "TagID", tagIDChanges)

		umIDChanges := SortByUniqueKey(&db.UserMark)
		UpdateIDs(db.BlockRange, "UserMarkID", umIDChanges)
		UpdateIDs(db.Note, "UserMarkID", umIDChanges)

		SortByUniqueKey(&db.BlockRange)

		noteIDChanges := SortByUniqueKey(&db.Note)
		UpdateIDs(db.TagMap, "NoteID", noteIDChanges)

		SortByUniqueKey(&db.TagMap)
	}

	// Check if all entries are equal.
	dbFields := reflect.ValueOf(dbCp).Elem()
	otherFields := reflect.ValueOf(otherCp).Elem()
	for i := 0; i < dbFields.NumField(); i++ {
		dbSlice := dbFields.Field(i)
		otherSlice := otherFields.Field(i)
		if !dbSlice.CanInterface() || !otherSlice.CanInterface() {
			continue
		}

		if dbSlice.Len() != otherSlice.Len() {
			fmt.Printf("Length of slices at index %d are not equal: %d vs %d\n", i, dbSlice.Len(), otherSlice.Len())
			return false
		}

		for j := 0; j < dbSlice.Len(); j++ {
			dElem := dbSlice.Index(j)
			oElem := otherSlice.Index(j)

			if dElem.IsNil() {
				if oElem.IsNil() {
					continue
				} else {
					return false
				}
			}

			if !dElem.MethodByName("Equals").Call([]reflect.Value{oElem})[0].Bool() {
				fmt.Println("Found different entries: ")
				left := spew.Sdump(dElem.Interface())
				right := spew.Sdump(oElem.Interface())
				fmt.Printf("%s \nvs\n %s", left, right)
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(left, right, true)
				fmt.Println("Diff:")
				fmt.Println(dmp.DiffPrettyText(diffs))
				return false
			}
		}
	}

	return true
}

// ImportJWLBackup unzips a given JW Library Backup file and imports the
// included SQLite DB to the Database struct
func (db *Database) ImportJWLBackup(filename string) error {
	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", "go-jwlm")
	if err != nil {
		return errors.Wrap(err, "Error while creating temporary directory")
	}
	defer os.RemoveAll(tmp)

	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		path := filepath.Join(tmp, file.Name)
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return errors.Wrap(err, "Error while copying files from backup to temporary folder")
		}
	}

	// Import manifest
	path := filepath.Join(tmp, manifestFilename)
	manifest := manifest{}
	if err := manifest.importManifest(path); err != nil {
		return errors.Wrap(err, "Error while importing manifest")
	}

	// Make sure that we support this backup version
	if err := manifest.validateManifest(); err != nil {
		return err
	}

	// Fill the Database with actual data
	path = filepath.Join(tmp, manifest.UserDataBackup.DatabaseName)
	return db.importSQLite(path)
}

// importSQLite imports a given SQLite DB into the Database struct
func (db *Database) importSQLite(filename string) error {
	// Open SQLite file as immutable to avoid locks (and therefore speed up import)
	sqlite, err := sql.Open("sqlite3", filename+"?immutable=1")
	if err != nil {
		return errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	// Make sure these tables are empty as we are not able to merge them yet.
	// Better to fail, than to risk losing data..
	emptyTables := []string{"PlaylistItem", "PlaylistItemChild", "PlaylistMedia"}
	for _, table := range emptyTables {
		count, err := getTableEntryCount(sqlite, table)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("Table %s is not empty. Merging of these entries are not supported yet", table)
		}
	}

	var wg sync.WaitGroup
	wg.Add(8)
	errorChan := make(chan error, 10)

	// Fill each table struct separately (did not find a DRYer solution yet..)
	go func() {
		mdl, err := fetchFromSQLite(sqlite, &BlockRange{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.BlockRange = BlockRange{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &Bookmark{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.Bookmark = Bookmark{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &InputField{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.InputField = InputField{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &Location{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.Location = Location{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &Note{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.Note = Note{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &Tag{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.Tag = Tag{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &TagMap{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.TagMap = TagMap{}.MakeSlice(mdl)
		wg.Done()
	}()

	go func() {
		mdl, err := fetchFromSQLite(sqlite, &UserMark{})
		if err != nil {
			errorChan <- err
			wg.Done()
			return
		}
		db.UserMark = UserMark{}.MakeSlice(mdl)
		wg.Done()
	}()

	wg.Wait()

	select {
	case err := <-errorChan:
		return err
	default:
		return nil
	}
}

// fetchFromSQLite fetches the entries for a given modelType and returns a slice
// of entries, for which the index corresponds to the ID in the SQLite DB
func fetchFromSQLite(sqlite *sql.DB, modelType Model) ([]Model, error) {
	// Create slice of correct size (number of entries)
	capacity, err := getSliceCapacity(sqlite, modelType)
	if err != nil {
		return nil, errors.Wrap(err, "Could not determine number of entries in SQLite database")
	}
	result := make([]Model, capacity)

	rows, err := sqlite.Query(fmt.Sprintf("SELECT * FROM %s", modelType.tableName()))
	if err != nil {
		return nil, errors.Wrap(err, "Error while querying SQLite database")
	}

	// Put entries in slice with the index coresponding to the ID in the SQLite DB
	i := 1
	defer rows.Close()
	for rows.Next() {
		var m Model
		switch tp := modelType.(type) {
		case *BlockRange:
			m = &BlockRange{}
		case *Bookmark:
			m = &Bookmark{}
		case *InputField:
			m = &InputField{pseudoID: i}
		case *Location:
			m = &Location{}
		case *Note:
			m = &Note{}
		case *Tag:
			m = &Tag{}
		case *TagMap:
			m = &TagMap{}
		case *UserMark:
			m = &UserMark{}
		default:
			panic(fmt.Sprintf("Fetching %T is not supported!", tp))
		}
		mn, err := m.scanRow(rows)
		if err != nil {
			return nil, errors.Wrap(err, "Error while scanning results from SQLite database")
		}
		result[mn.ID()] = mn
		i++
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error while scanning results from SQLite database")
	}

	return result, nil
}

// getTableEntryCount returns the number of entries in a given table
func getTableEntryCount(sqlite *sql.DB, tableName string) (int, error) {
	var count int
	err := sqlite.QueryRow(fmt.Sprintf("SELECT Count(*) FROM %s", tableName)).Scan(&count)
	if err != nil {
		return 0, errors.Wrapf(err, "Error while determing entry count of table %s", tableName)
	}

	return count, nil
}

// getSliceCapacity determines the needed capacity for a slice from a table
// by looking at the highest ID in the DB. If the table does not have a ID
// column, it will simply count the number of entries.
func getSliceCapacity(sqlite *sql.DB, modelType Model) (int, error) {
	var query string
	if modelType.idName() != "" {
		query = fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC LIMIT 1",
			modelType.idName(), modelType.tableName(), modelType.idName())
	} else {
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s", modelType.tableName())
	}
	row, err := sqlite.Query(query)
	if err != nil {
		return 0, err
	}
	defer row.Close()

	capacity := 0
	for row.Next() {
		if err := row.Scan(&capacity); err != nil {
			return 0, err
		}
	}

	// Index in DB starts with 1, so 0 is always nil
	return capacity + 1, nil
}

// ExportJWLBackup creates a .jwlibrary backup file out of a Database{} struct
func (db *Database) ExportJWLBackup(filename string) error {
	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", "go-jwlm")
	if err != nil {
		return errors.Wrap(err, "Error while creating temporary directory")
	}
	defer os.RemoveAll(tmp)

	// Create user_data.db
	dbPath := filepath.Join(tmp, "user_data.db")
	if err := db.saveToNewSQLite(dbPath); err != nil {
		return errors.Wrap(err, "Could not create SQLite database for exporting")
	}

	// Create manifest.json
	manifestPath := filepath.Join(tmp, manifestFilename)
	mfst, err := generateManifest("go-jwlm", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error while generating manifest")
	}
	if err := mfst.exportManifest(manifestPath); err != nil {
		return errors.Wrap(err, "Error while creating manifest.json")
	}

	// Store files in .jwlibrary (zip)-file
	files := []string{dbPath, manifestPath}
	if err := zipFiles(filename, files); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while storing files in zip archive %s", filename))
	}

	return nil
}

// SaveToNewSQLite creates a new SQLite database with the JW Library scheme
// and saves all entries of the Database{} struct to it
func (db *Database) saveToNewSQLite(filename string) error {
	if err := createEmptySQLiteDB(filename); err != nil {
		return errors.Wrap(err, "Error while creating new empty SQLite database")
	}

	sqlite, err := sql.Open("sqlite3", filename)
	if err != nil {
		return errors.Wrap(err, "Error while opening SQLite database")
	}
	defer sqlite.Close()

	// For every field of the Database{} struct, create a []model slice
	// and use it to insert its entries to the new SQLite DB
	dbFields := reflect.ValueOf(db).Elem()
	for j := 0; j < dbFields.NumField(); j++ {
		slice := dbFields.Field(j).Interface()
		mdl, err := MakeModelSlice(slice)
		if err != nil {
			return err
		}
		if err := insertEntries(sqlite, mdl); err != nil {
			return errors.Wrapf(err, "Error while inserting entries of field %d", j)
		}
	}

	// Update LastModified
	lastModified := time.Now().Format("2006-01-02T15:04:05-07:00")
	_, err = sqlite.Exec(fmt.Sprintf("UPDATE LastModified SET LastModified = \"%s\" WHERE LastModified = (SELECT * FROM LastModified)", lastModified))
	if err != nil {
		return errors.Wrap(err, "Error while updating LastModified")
	}

	// Vacuum to clean up SQLite DB
	_, err = sqlite.Exec("VACUUM")
	if err != nil {
		return errors.Wrap(err, "Error while vacuuming SQLite DB")
	}

	return nil
}

// insertEntries INSERTs entries of []model into a given SQLite database.
// It does it by dynamically parsing all fields of a struct implementing
// model using reflection and creating a query for SQLite out of it.
func insertEntries(sqlite *sql.DB, m []Model) error {
	if len(m) == 0 {
		return nil
	}

	// Figure out tableName and rowCount. As we don't know for sure,
	// which entry will be nil-pointer, we just try until we find a non-empty
	// one and call the functions there.
	tableName := ""
	rowCount := 0
	foundEntry := false
	for _, mdl := range m {
		if reflect.ValueOf(mdl).Elem().IsValid() {
			tableName = mdl.tableName()
			foundEntry = true

			// Count number of fields that don't have the "ignore" tag set
			reflTypes := reflect.TypeOf(mdl).Elem()
			for j := 0; j < reflTypes.NumField(); j++ {
				if _, ignore := reflTypes.Field(j).Tag.Lookup("ignore"); ignore {
					continue
				}
				rowCount++
			}
			break
		}
	}
	// If slice is empty, we don't need to continue
	if !foundEntry {
		return nil
	}

	tx, err := sqlite.Begin()
	if err != nil {
		return err
	}

	// Dynamically add all column-names of the struct to the query
	query := fmt.Sprintf("INSERT INTO %s VALUES (", tableName)
	for i := 0; i < rowCount; i++ {
		query += "?"
		if i+1 < rowCount {
			query += ", "
		}
	}
	query += ")"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.Wrapf(err, "Error while preparing query %s", query)
	}
	defer stmt.Close()

	for _, entry := range m {
		// Prepare struct for ingestion with stmt.Exec
		values := make([]interface{}, rowCount)
		reflValues := reflect.ValueOf(entry).Elem()
		reflTypes := reflect.TypeOf(entry).Elem()

		// Check if entry is actually a nil-pointer and shouldn't be considered
		if !reflValues.IsValid() {
			continue
		}

		// Add all fields of the struct to the values slice, so we can ingest them later
		for j := 0; j < reflValues.NumField(); j++ {
			// If struct field has `ignore` tag then skip it
			if _, ignore := reflTypes.Field(j).Tag.Lookup("ignore"); ignore {
				continue
			}
			v := reflValues.Field(j).Interface()
			values[j] = v
		}

		if _, err := stmt.Exec(values...); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Could not insert entry %v", entry))
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Error while commiting entries")
	}

	return nil
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
