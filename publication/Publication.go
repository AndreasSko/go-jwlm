package publication

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AndreasSko/go-jwlm/model"
	snippets "github.com/AndreasSko/jwpub-snippets"
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

// GetSnippet returns the snippet of the UserMarkBlockRange in the given publication.
// It uses logic in a privat repository, which only the GitHub Action has access to.
// In other cases the fake-implementation is used which by default just returns an error.
func GetSnippet(publPath string, publ Publication, loc model.Location, umbr model.UserMarkBlockRange) ([]string, error) {
	query := snippets.SnippetQuery{
		Publication: snippets.Publication{
			ID:                    publ.ID,
			PublicationRootKeyID:  publ.PublicationRootKeyID,
			MepsLanguageID:        publ.MepsLanguageID,
			PublicationTypeID:     publ.PublicationTypeID,
			IssueTagNumber:        publ.IssueTagNumber,
			Title:                 publ.Title,
			IssueTitle:            publ.IssueTitle,
			ShortTitle:            publ.ShortTitle,
			CoverTitle:            publ.CoverTitle,
			UndatedTitle:          publ.UndatedTitle,
			UndatedReferenceTitle: publ.UndatedReferenceTitle,
			Year:                  publ.Year,
			Symbol:                publ.Symbol,
			KeySymbol:             publ.KeySymbol,
			Reserved:              publ.Reserved,
		},
		Location:           loc,
		UserMarkBlockRange: umbr,
	}

	return snippets.GetSnippet(publPath, query)
}

// GetPublicationPath generates the filename of the publication and checks if it
// exists in the publDir
func GetPublicationPath(publ Publication, publDir string) (string, error) {
	language, err := lookupMepsLanguage(publ.MepsLanguageID)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s_%s_%d", publ.KeySymbol.String, language.Symbol, publ.IssueTagNumber)
	for strings.HasSuffix(filename, "0") {
		filename = strings.TrimSuffix(filename, "0")
	}
	filename = strings.TrimSuffix(filename, "_")
	filename += ".db"

	path := filepath.Join(publDir, filename)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("did not find publication: %w", err)
	}

	return path, nil
}

// MarshalJSON returns the JSON encoding of the entry
func (m Publication) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID                    int    `json:"id"`
		PublicationRootKeyID  int    `json:"publicationRootKeyId"`
		MepsLanguageID        int    `json:"mepsLanguageId"`
		PublicationTypeID     int    `json:"publicationTypeId"`
		IssueTagNumber        int    `json:"issueTagNumber"`
		Title                 string `json:"title"`
		IssueTitle            string `json:"issueTitle"`
		ShortTitle            string `json:"shortTitle"`
		CoverTitle            string `json:"coverTitle"`
		UndatedTitle          string `json:"undatedTitle"`
		UndatedReferenceTitle string `json:"undatedReferenceTitle"`
		Year                  int    `json:"year"`
		Symbol                string `json:"symbol"`
		KeySymbol             string `json:"keySymbol"`
		Reserved              int    `json:"reserved"`
	}{
		ID:                    m.ID,
		PublicationRootKeyID:  m.PublicationRootKeyID,
		MepsLanguageID:        m.MepsLanguageID,
		PublicationTypeID:     m.PublicationTypeID,
		IssueTagNumber:        m.IssueTagNumber,
		Title:                 m.Title,
		IssueTitle:            m.IssueTitle.String,
		ShortTitle:            m.ShortTitle,
		CoverTitle:            m.CoverTitle.String,
		UndatedTitle:          m.UndatedTitle.String,
		UndatedReferenceTitle: m.UndatedReferenceTitle.String,
		Year:                  m.Year,
		Symbol:                m.Symbol,
		KeySymbol:             m.KeySymbol.String,
		Reserved:              m.Reserved,
	})
}

func (m *Publication) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID                    int    `json:"id"`
		PublicationRootKeyID  int    `json:"publicationRootKeyId"`
		MepsLanguageID        int    `json:"mepsLanguageId"`
		PublicationTypeID     int    `json:"publicationTypeId"`
		IssueTagNumber        int    `json:"issueTagNumber"`
		Title                 string `json:"title"`
		IssueTitle            string `json:"issueTitle,,omitempty"`
		ShortTitle            string `json:"shortTitle"`
		CoverTitle            string `json:"coverTitle,omitempty"`
		UndatedTitle          string `json:"undatedTitle,omitempty"`
		UndatedReferenceTitle string `json:"undatedReferenceTitle,omitempty"`
		Year                  int    `json:"year"`
		Symbol                string `json:"symbol"`
		KeySymbol             string `json:"keySymbol,omitempty"`
		Reserved              int    `json:"reserved"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	m.ID = aux.ID
	m.PublicationRootKeyID = aux.PublicationRootKeyID
	m.MepsLanguageID = aux.MepsLanguageID
	m.PublicationTypeID = aux.PublicationTypeID
	m.IssueTagNumber = aux.IssueTagNumber
	m.Title = aux.Title
	m.IssueTitle = sql.NullString{String: aux.IssueTitle, Valid: aux.IssueTitle != ""}
	m.ShortTitle = aux.ShortTitle
	m.CoverTitle = sql.NullString{String: aux.CoverTitle, Valid: aux.CoverTitle != ""}
	m.UndatedTitle = sql.NullString{String: aux.UndatedTitle, Valid: aux.UndatedTitle != ""}
	m.UndatedReferenceTitle = sql.NullString{String: aux.UndatedReferenceTitle, Valid: aux.UndatedReferenceTitle != ""}
	m.Year = aux.Year
	m.Symbol = aux.Symbol
	m.KeySymbol = sql.NullString{String: aux.KeySymbol, Valid: aux.KeySymbol != ""}
	m.Reserved = aux.Reserved

	return nil
}
