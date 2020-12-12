package publication

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"

	// Register SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// Publication represents a publication with all
// its information from the catalogDB
type Publication struct {
	ID                    int
	PublicationRootKeyID  int
	MepsLanguageID        int
	PublicationTypeID     int
	IssueTagNumber        int
	Title                 string
	IssueTitle            sql.NullString
	ShortTitle            string
	CoverTitle            sql.NullString
	UndatedTitle          sql.NullString
	UndatedReferenceTitle sql.NullString
	Year                  int
	Symbol                string
	KeySymbol             sql.NullString
	Reserved              int
}

// Lookup represents a lookup for a publication.
// This query can contain various fields.
type Lookup struct {
	DocumentID     int
	KeySymbol      string
	IssueTagNumber int
	MepsLanguage   int
}

// LookupPublication looks up a publication from catalogDB located at dbPath
func LookupPublication(dbPath string, query Lookup) (Publication, error) {
	// Check if file exists
	if _, err := os.Stat(dbPath); err != nil {
		return Publication{}, fmt.Errorf("CatalogDB does not exist at %s", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath+"?immutable=1")
	if err != nil {
		return Publication{}, errors.Wrap(err, "Error while opening SQLite database")
	}
	defer db.Close()

	return lookupPublication(db, query)
}

func lookupPublication(db *sql.DB, query Lookup) (Publication, error) {
	var row *sql.Row
	if query.DocumentID != 0 {
		stmt, err := db.Prepare("SELECT P.* " +
			"FROM Publication AS P, PublicationDocument AS PD " +
			"WHERE P.Id = PD.PublicationId AND PD.DocumentId = ? AND P.MepsLanguageId = ?")
		if err != nil {
			return Publication{}, errors.Wrap(err, "Error while preparing query")
		}
		row = stmt.QueryRow(query.DocumentID, query.MepsLanguage)
	} else {
		stmt, err := db.Prepare("SELECT * FROM Publication WHERE KeySymbol = ? AND MepsLanguageId = ? AND IssueTagNumber = ?")
		if err != nil {
			return Publication{}, errors.Wrap(err, "Error while preparing query")
		}
		row = stmt.QueryRow(query.KeySymbol, query.MepsLanguage, query.IssueTagNumber)
	}

	publ := Publication{}
	err := row.Scan(&publ.PublicationRootKeyID,
		&publ.MepsLanguageID,
		&publ.PublicationTypeID,
		&publ.IssueTagNumber,
		&publ.Title,
		&publ.IssueTitle,
		&publ.ShortTitle,
		&publ.CoverTitle,
		&publ.UndatedTitle,
		&publ.UndatedReferenceTitle,
		&publ.Year,
		&publ.Symbol,
		&publ.KeySymbol,
		&publ.Reserved,
		&publ.ID)
	if err != nil {
		return Publication{}, errors.Wrap(err, "Error while scanning row for publication")
	}

	return publ, nil
}
