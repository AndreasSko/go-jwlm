package publication

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var publicationColumns = []string{
	"PublicationRootKeyId",
	"MepsLanguageId",
	"PublicationTypeId",
	"IssueTagNumber",
	"Title",
	"IssueTitle",
	"ShortTitle",
	"CoverTitle",
	"UndatedTitle",
	"UndatedReferenceTitle",
	"Year",
	"Symbol",
	"KeySymbol",
	"Reserved",
	"Id",
}

func TestLookupPublication(t *testing.T) {
	var tests = []struct {
		input       Lookup
		expected    Publication
		expectError bool
	}{
		{
			input: Lookup{
				KeySymbol:    "cl",
				MepsLanguage: 0,
			},
			expected: Publication{
				ID:                    67,
				PublicationRootKeyID:  64,
				MepsLanguageID:        0,
				PublicationTypeID:     2,
				IssueTagNumber:        0,
				Title:                 "Draw Close to Jehovah",
				ShortTitle:            "Close to Jehovah",
				UndatedReferenceTitle: sql.NullString{"Close to Jehovah", true},
				Year:                  2014,
				Symbol:                "cl",
				KeySymbol:             sql.NullString{"cl", true},
				Reserved:              0,
			},
		},
		{
			input: Lookup{
				DocumentID:   1102002020,
				MepsLanguage: 0,
			},
			expected: Publication{
				ID:                    67,
				PublicationRootKeyID:  64,
				MepsLanguageID:        0,
				PublicationTypeID:     2,
				IssueTagNumber:        0,
				Title:                 "Draw Close to Jehovah",
				ShortTitle:            "Close to Jehovah",
				UndatedReferenceTitle: sql.NullString{"Close to Jehovah", true},
				Year:                  2014,
				Symbol:                "cl",
				KeySymbol:             sql.NullString{"cl", true},
				Reserved:              0,
			},
		},
		{
			input: Lookup{
				DocumentID:   1102002020,
				MepsLanguage: 1,
			},
			expected:    Publication{},
			expectError: true,
		},
		{
			input: Lookup{
				KeySymbol:    "cl",
				MepsLanguage: 1,
			},
			expected: Publication{
				ID:                    129,
				PublicationRootKeyID:  64,
				MepsLanguageID:        1,
				PublicationTypeID:     2,
				IssueTagNumber:        0,
				Title:                 "Acerquémonos a Jehová",
				ShortTitle:            "Acerquémonos a Jehová",
				UndatedReferenceTitle: sql.NullString{"Acerquémonos a Jehová", true},
				Year:                  2014,
				Symbol:                "cl",
				KeySymbol:             sql.NullString{"cl", true},
				Reserved:              0,
			},
		},
		{
			input: Lookup{
				IssueTagNumber: 20210200,
				KeySymbol:      "w",
				MepsLanguage:   0,
			},
			expected: Publication{
				ID:                    305097,
				PublicationRootKeyID:  780,
				MepsLanguageID:        0,
				PublicationTypeID:     14,
				IssueTagNumber:        20210200,
				Title:                 "The Watchtower Announcing Jehovah’s Kingdom (Study)—2021",
				IssueTitle:            sql.NullString{"The Watchtower, February 2021", true},
				ShortTitle:            "The Watchtower (Study) (2021)",
				CoverTitle:            sql.NullString{"Study Articles for April 5 to May 2", true},
				UndatedTitle:          sql.NullString{"The Watchtower—Study Edition", true},
				UndatedReferenceTitle: sql.NullString{"The Watchtower (Study)", true},
				Year:                  2021,
				Symbol:                "w21",
				KeySymbol:             sql.NullString{"w", true},
				Reserved:              0,
			},
		},
	}

	path := filepath.Join("testdata", "catalog.db")
	for _, test := range tests {
		res, err := LookupPublication(path, test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expected, res)
	}
}

func Test_lookupPublication(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	publication := Publication{
		PublicationRootKeyID:  1,
		MepsLanguageID:        2,
		PublicationTypeID:     3,
		IssueTagNumber:        0,
		Title:                 "Title",
		IssueTitle:            sql.NullString{"IssueTitle", true},
		ShortTitle:            "ShortTitle",
		CoverTitle:            sql.NullString{"CoverTitle", true},
		UndatedTitle:          sql.NullString{"UndatedTitle", true},
		UndatedReferenceTitle: sql.NullString{"UndatedReferenceTitle", true},
		Year:                  2020,
		Symbol:                "Symbol",
		KeySymbol:             sql.NullString{"KeySymbol", true},
		Reserved:              0,
		ID:                    1,
	}

	mock.ExpectPrepare("SELECT P.\\* FROM Publication AS P, PublicationDocument AS PD " +
		"WHERE P.Id = PD.PublicationId AND PD.DocumentId = \\? AND P.MepsLanguageId = \\?").ExpectQuery().
		WillReturnRows(mock.NewRows(publicationColumns).AddRow(publication.PublicationRootKeyID,
			publication.MepsLanguageID,
			publication.PublicationTypeID,
			publication.IssueTagNumber,
			publication.Title,
			publication.IssueTitle,
			publication.ShortTitle,
			publication.CoverTitle,
			publication.UndatedTitle,
			publication.UndatedReferenceTitle,
			publication.Year,
			publication.Symbol,
			publication.KeySymbol,
			publication.Reserved,
			publication.ID))

	res, err := lookupPublication(db, Lookup{DocumentID: 1})
	assert.NoError(t, err)
	assert.Equal(t, publication, res)

	mock.ExpectPrepare("SELECT \\* FROM Publication WHERE KeySymbol = \\? AND MepsLanguageId = \\? AND IssueTagNumber = \\?").ExpectQuery().
		WillReturnRows(mock.NewRows(publicationColumns).AddRow(publication.PublicationRootKeyID,
			publication.MepsLanguageID,
			publication.PublicationTypeID,
			publication.IssueTagNumber,
			publication.Title,
			publication.IssueTitle,
			publication.ShortTitle,
			publication.CoverTitle,
			publication.UndatedTitle,
			publication.UndatedReferenceTitle,
			publication.Year,
			publication.Symbol,
			publication.KeySymbol,
			publication.Reserved,
			publication.ID))
	res, err = lookupPublication(db, Lookup{KeySymbol: "KeySymbol"})
	assert.NoError(t, err)
	assert.Equal(t, publication, res)
}

func TestPublication_MarshalJSON(t *testing.T) {
	publ := Publication{
		ID:                    1,
		PublicationRootKeyID:  2,
		MepsLanguageID:        3,
		PublicationTypeID:     4,
		IssueTagNumber:        5,
		Title:                 "6",
		IssueTitle:            sql.NullString{},
		ShortTitle:            "8",
		CoverTitle:            sql.NullString{"9", true},
		UndatedTitle:          sql.NullString{"10", true},
		UndatedReferenceTitle: sql.NullString{"11", true},
		Year:                  12,
		Symbol:                "13",
		KeySymbol:             sql.NullString{"14", true},
		Reserved:              15,
	}

	expected := `{"id":1,"publicationRootKeyId":2,"mepsLanguageId":3,"publicationTypeId":4,"issueTagNumber":5,"title":"6","issueTitle":"","shortTitle":"8","coverTitle":"9","undatedTitle":"10","undatedReferenceTitle":"11","year":12,"symbol":"13","keySymbol":"14","reserved":15}`

	jsn, err := json.Marshal(publ)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(jsn))
}
