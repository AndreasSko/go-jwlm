package model

import (
	"database/sql"
	"runtime"
	"sync"
)

var once sync.Once
var persistence Persistence

type Persistence interface {
	CreateTempStorage(prefix string) (path string, err error)
	StoreSQLiteDB(filename string, dbData []byte) (fullFileName string, err error)
	OpenSQLiteDB(fullFileName string) (*sql.DB, error)
	RetrieveSQLiteData(fullFileName string) ([]byte, error)
	StoreJWLBackup(fullFileName string, archiveData []byte) error
	ProcessJWLBackup(fullFileName string, exportPath string) error
	GetFile(fullFileName string) (filename string, data []byte, err error)
	WriteFile(fullFileName string, data []byte) error
	CleanupPath(path string) error
}

func GetPersistence() Persistence {
	once.Do(func() {
		if runtime.GOOS == "js" {
			persistence = getJsPersistence()
		} else {
			persistence = getFsPersistence()
		}
	})

	return persistence

}
