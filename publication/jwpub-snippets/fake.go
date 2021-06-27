package snippets

import (
	"database/sql"
	"fmt"

	"github.com/AndreasSko/go-jwlm/model"
)

const FakeImplementation = true

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
type SnippetQuery struct {
	Publication        Publication
	Location           model.Location
	UserMarkBlockRange model.UserMarkBlockRange
}

// GetSnippet is a fake implementation and only used as a general replacement
// of github.com/AndreasSko/jwpub-snippets so go-jwlm is able to be built with
// or without it
func GetSnippet(publPath string, query SnippetQuery) ([]string, error) {
	return nil, fmt.Errorf("Using dummy implementation")
}
