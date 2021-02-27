// +build js

package model

import (
	// Register SQLite driver
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall/js"

	_ "github.com/fritzbauer/go-sqlite3-js"
	"github.com/pkg/errors"
)

const jsSqliteVarPrefix = "_go-jwlm_db"

type jsPersistence struct {
	storage map[string]*PersistedFolder
}

type PersistedFolder struct {
	Files map[string]*PersistedFile
}

type PersistedFile struct {
	Name string
	Data []byte
}

func getJsPersistence() Persistence {
	pers := jsPersistence{storage: make(map[string]*PersistedFolder)}
	return &pers
}

func getFsPersistence() Persistence {
	panic("getFsPersistence call in js runtime")
}

func (pers *jsPersistence) CreateTempStorage(prefix string) (path string, err error) {
	//find the first unused foldername
	var i int
	for i = 0; pers.storage[fmt.Sprintf("%s_%d", prefix, i)] != nil; i++ {
	}
	path = fmt.Sprintf("%s_%d", prefix, i)
	pers.storage[path] = &PersistedFolder{Files: make(map[string]*PersistedFile)}
	return path, nil
}

func (pers *jsPersistence) StoreSQLiteDB(filename string, dbData []byte) (fullFileName string, err error) {
	//TODO think of proper errorhandling
	jsName := fmt.Sprintf("%s_%d_%s", jsSqliteVarPrefix, 0, filename)

	//prevent using an already used JS variable
	for i := 1; js.Global().Get(jsName).Truthy(); i++ {
		jsName = fmt.Sprintf("%s_%d_%s", jsSqliteVarPrefix, i, filename)
	}

	arr := js.Global().Get("Uint8Array").New(len(dbData))
	js.CopyBytesToJS(arr, dbData)
	js.Global().Set(jsName, arr)
	return jsName, nil
}

func (pers *jsPersistence) OpenSQLiteDB(fullFileName string) (*sql.DB, error) {
	if strings.Contains(fullFileName, "?") { //remove ?immutable=1
		fullFileName = strings.Split(fullFileName, "?")[0]
	}
	_, data, err := pers.GetFile(fullFileName)
	jsStorageVariableName, err := pers.StoreSQLiteDB(fullFileName, data)
	if err != nil {
		return nil, errors.Wrap(err, "Could not store SQLite db")
	}

	//add a "file" to store the mapping of filepath to js var name
	err = pers.WriteFile(fullFileName, []byte(jsStorageVariableName))
	if err != nil {
		return nil, errors.Wrap(err, "Error storing jsStorageVariableName")
	}

	return sql.Open("sqlite3_js", jsStorageVariableName)
}

func (pers *jsPersistence) RetrieveSQLiteData(jsVarName string) ([]byte, error) {
	//TODO think of proper errorhandling
	dbMap := js.Global().Get("_go_sqlite_dbs")
	jsData := dbMap.Call("get", jsVarName).Call("export")
	data := make([]byte, jsData.Get("byteLength").Int())
	js.CopyBytesToGo(data, jsData)
	return data, nil
}

func (pers *jsPersistence) StoreJWLBackup(fullFileName string, archiveData []byte) error {
	path, fileName := evaluateFullFileName(fullFileName)

	folder, ok := pers.storage[path]
	if !ok {
		pers.storage[path] = &PersistedFolder{Files: make(map[string]*PersistedFile)}
		folder = pers.storage[path]
	}

	folder.Files[fileName] = &PersistedFile{Name: fileName, Data: archiveData}
	//pers.printStorage() //debug
	return nil
}

func (pers *jsPersistence) ProcessJWLBackup(fullFileName string, exportPath string) error {

	_, data, err := pers.GetFile(fullFileName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not find JW Library backup at %v", fullFileName))
	}
	reader := bytes.NewReader(data)

	r, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return errors.Wrap(err, "Could not read zip")
	}

	for _, file := range r.File {
		fileReader, err := file.Open()
		if err != nil {
			return errors.Wrap(err, "Error while opening zip file")
		}
		defer fileReader.Close()

		var buf bytes.Buffer
		_, err = io.Copy(bufio.NewWriter(&buf), fileReader)
		if err != nil {
			return errors.Wrap(err, "Error while storing files from backup ")
		}

		path := filepath.Join(exportPath, file.Name)
		pers.WriteFile(path, buf.Bytes())
	}
	//pers.printStorage() //debug
	return nil
}

func (pers *jsPersistence) GetFile(fullFileName string) (filename string, data []byte, err error) {
	path, fileName := evaluateFullFileName(fullFileName)
	file := pers.storage[path].Files[fileName]
	if file == nil {
		return "", nil, errors.Errorf("Could not find file '%v' at %v", fileName, path)
	}

	if strings.HasPrefix(string(file.Data), jsSqliteVarPrefix) {
		mergedData, err := pers.RetrieveSQLiteData(string(file.Data))
		if err != nil {
			return "", nil, errors.Wrap(err, "Error while getting SQLite DB")
		}
		return fileName, mergedData, nil
	}

	//fmt.Printf("Returning %s; Length: %d\n", file.Name, len(file.Data)) //debug
	return file.Name, file.Data, nil
}

func (pers *jsPersistence) WriteFile(fullFileName string, data []byte) error {
	path, fileName := evaluateFullFileName(fullFileName)

	folder := pers.getFolder(path)

	folder.Files[fileName] = &PersistedFile{Name: fileName, Data: data} //this is silently overwriting any existing file
	return nil
}

func (pers *jsPersistence) CleanupPath(path string) error {
	delete(pers.storage, path)
	return nil
}

func (pers *jsPersistence) getFolder(path string) *PersistedFolder {
	folder, ok := pers.storage[path]
	if !ok {
		pers.storage[path] = &PersistedFolder{Files: make(map[string]*PersistedFile)}
		folder = pers.storage[path]
	}
	return folder
}

func evaluateFullFileName(fullFileName string) (path string, fileName string) {
	fullfileNameParts := strings.Split(fullFileName, string(os.PathSeparator))
	path = strings.Join(fullfileNameParts[:len(fullfileNameParts)-1], string(os.PathSeparator))
	fileName = fullfileNameParts[len(fullfileNameParts)-1]
	//fmt.Printf("Splitted '%s' into '%s' and '%s'\n", fullFileName, path, fileName) //debug
	return path, fileName
}

//Debugging print
func (pers *jsPersistence) printStorage() {
	for folderName, folder := range pers.storage {
		for filename, file := range folder.Files {
			fmt.Printf("PrintStorage: %s/%s File.Name: %s; Length: %d\n", folderName, filename, file.Name, len(file.Data))
		}
	}
}
