// +build !windows

package gomobile

import (
	"path/filepath"
	"testing"

	"github.com/tj/assert"
)

func TestLookupPublication(t *testing.T) {
	var tests = []struct {
		input    *PublicationLookup
		expected string
	}{
		{
			input: &PublicationLookup{
				KeySymbol:    "cl",
				MepsLanguage: 0,
			},
			expected: `{"id":67,"publicationRootKeyId":64,"mepsLanguageId":0,"publicationTypeId":2,"issueTagNumber":0,"title":"Draw Close to Jehovah","issueTitle":"","shortTitle":"Close to Jehovah","coverTitle":"","undatedTitle":"","undatedReferenceTitle":"Close to Jehovah","year":2014,"symbol":"cl","keySymbol":"cl","reserved":0}`,
		},
		{
			input: &PublicationLookup{
				DocumentID:   1102002020,
				MepsLanguage: 0,
			},
			expected: `{"id":67,"publicationRootKeyId":64,"mepsLanguageId":0,"publicationTypeId":2,"issueTagNumber":0,"title":"Draw Close to Jehovah","issueTitle":"","shortTitle":"Close to Jehovah","coverTitle":"","undatedTitle":"","undatedReferenceTitle":"Close to Jehovah","year":2014,"symbol":"cl","keySymbol":"cl","reserved":0}`,
		},
		{
			input: &PublicationLookup{
				DocumentID:   1102002020,
				MepsLanguage: 1,
			},
			expected: "",
		},
		{
			input: &PublicationLookup{
				KeySymbol:    "cl",
				MepsLanguage: 1,
			},
			expected: `{"id":129,"publicationRootKeyId":64,"mepsLanguageId":1,"publicationTypeId":2,"issueTagNumber":0,"title":"Acerquémonos a Jehová","issueTitle":"","shortTitle":"Acerquémonos a Jehová","coverTitle":"","undatedTitle":"","undatedReferenceTitle":"Acerquémonos a Jehová","year":2014,"symbol":"cl","keySymbol":"cl","reserved":0}`,
		},
		{
			input: &PublicationLookup{
				IssueTagNumber: 20210200,
				KeySymbol:      "w",
				MepsLanguage:   0,
			},
			expected: `{"id":305097,"publicationRootKeyId":780,"mepsLanguageId":0,"publicationTypeId":14,"issueTagNumber":20210200,"title":"The Watchtower Announcing Jehovah’s Kingdom (Study)—2021","issueTitle":"The Watchtower, February 2021","shortTitle":"The Watchtower (Study) (2021)","coverTitle":"Study Articles for April 5 to May 2","undatedTitle":"The Watchtower—Study Edition","undatedReferenceTitle":"The Watchtower (Study)","year":2021,"symbol":"w21","keySymbol":"w","reserved":0}`,
		},
	}

	path := filepath.Join("../publication/testdata", "catalog.db")
	for _, test := range tests {
		res := LookupPublication(path, test.input)
		assert.Equal(t, test.expected, res)
	}
}
