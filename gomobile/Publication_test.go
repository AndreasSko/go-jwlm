// +build !windows

package gomobile

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	snippets "github.com/AndreasSko/jwpub-snippets"
	"github.com/stretchr/testify/assert"
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

func TestMergeConflictsWrapper_GetSnippet(t *testing.T) {
	if snippets.FakeImplementation {
		fmt.Println("Fake implementation used. Skipping test")
		return
	}

	exampleDB := &model.Database{
		BlockRange: []*model.BlockRange{
			nil,
			{
				BlockRangeID: 1,
				Identifier:   14,
				StartToken:   sql.NullInt32{Int32: 0, Valid: true},
				EndToken:     sql.NullInt32{Int32: 1, Valid: true},
				UserMarkID:   1,
			},
			{
				BlockRangeID: 1,
				Identifier:   15,
				StartToken:   sql.NullInt32{Int32: 0, Valid: true},
				EndToken:     sql.NullInt32{Int32: 1, Valid: true},
				UserMarkID:   1,
			},
		},
		Location: []*model.Location{
			nil,
			{
				LocationID:   1,
				DocumentID:   sql.NullInt32{Int32: 1102020991, Valid: true},
				KeySymbol:    sql.NullString{String: "syr20", Valid: true},
				MepsLanguage: 0,
				LocationType: 0,
			},
		},
		UserMark: []*model.UserMark{
			nil,
			{
				UserMarkID: 1,
				LocationID: 1,
			},
		},
	}
	type fields struct {
		DBWrapper         *DatabaseWrapper
		conflicts         map[string]merger.MergeConflict
		unsolvedConflicts map[string]bool
		solutions         map[string]merger.MergeSolution
	}
	type args struct {
		catalogDir  string
		publDir     string
		conflictKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get 2020 Service Year Report",
			fields: fields{
				DBWrapper: &DatabaseWrapper{
					left: exampleDB,
				},
				conflicts: map[string]merger.MergeConflict{
					"syr20": {
						Left: &model.UserMarkBlockRange{
							UserMark: exampleDB.UserMark[1],
							BlockRanges: []*model.BlockRange{
								exampleDB.BlockRange[1],
								exampleDB.BlockRange[2],
							},
						},
					},
				},
			},
			args: args{
				catalogDir:  "./testdata/catalog.db",
				publDir:     "./testdata/syr20_E.db",
				conflictKey: "syr20",
			},
			want: "[\"Average Bible Studies Each Month:  7,705,765\",\"During the 2020 service year, Jehovah’s Witnesses spent $231\u00a0million in caring for special pioneers, missionaries, and circuit overseers in their field service assignments. Worldwide, a total of 20,994 ordained ministers staff the branch facilities. All are members of the Worldwide Order of Special Full-Time Servants of Jehovah’s Witnesses.\"]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcw := &MergeConflictsWrapper{
				DBWrapper:         tt.fields.DBWrapper,
				conflicts:         tt.fields.conflicts,
				unsolvedConflicts: tt.fields.unsolvedConflicts,
				solutions:         tt.fields.solutions,
			}
			got, err := mcw.GetSnippet(tt.args.catalogDir, tt.args.publDir, tt.args.conflictKey)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPublicationPath(t *testing.T) {
	type args struct {
		publJSON string
		publDir  string
	}
	tests := []struct {
		name         string
		args         args
		wantContains string
		wantErr      bool
	}{
		{
			name: "WT Chinese - exists",
			args: args{
				publJSON: `{"id":305097,"publicationRootKeyId":780,"mepsLanguageId":43,"publicationTypeId":14,"issueTagNumber":20210200,"title":"The Watchtower Announcing Jehovah’s Kingdom (Study)—2021","issueTitle":"The Watchtower, February 2021","shortTitle":"The Watchtower (Study) (2021)","coverTitle":"Study Articles for April 5 to May 2","undatedTitle":"The Watchtower—Study Edition","undatedReferenceTitle":"The Watchtower (Study)","year":2021,"symbol":"w21","keySymbol":"w","reserved":0}`,
				publDir:  "./testdata",
			},
			wantContains: "testdata/w_CH_202102.db",
		},
		{
			name: "WT German - does not exist",
			args: args{
				publJSON: `{"id":305097,"publicationRootKeyId":780,"mepsLanguageId":2,"publicationTypeId":14,"issueTagNumber":20210200,"title":"The Watchtower Announcing Jehovah’s Kingdom (Study)—2021","issueTitle":"The Watchtower, February 2021","shortTitle":"The Watchtower (Study) (2021)","coverTitle":"Study Articles for April 5 to May 2","undatedTitle":"The Watchtower—Study Edition","undatedReferenceTitle":"The Watchtower (Study)","year":2021,"symbol":"w21","keySymbol":"w","reserved":0}`,
				publDir:  "./testdata",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPublicationPath(tt.args.publJSON, tt.args.publDir)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Contains(t, got, tt.wantContains)
		})
	}
}

func TestDownloadPublication(t *testing.T) {
	type args struct {
		publJSON string
	}
	tests := []struct {
		name        string
		args        args
		filename    string
		minFileSize int
		wantErr     bool
	}{
		{
			name: "Download 2020 service year report",
			args: args{
				publJSON: `{"id":317733,"publicationRootKeyId":786,"mepsLanguageId":0,"publicationTypeId":2,"issueTagNumber":0,"title":"2020 Service Year Report of Jehovah’s Witnesses Worldwide","issueTitle":"","shortTitle":"2020 Service Year Report","coverTitle":"","undatedTitle":"","undatedReferenceTitle":"2020 Service Year Report","year":2020,"symbol":"syr20","keySymbol":"syr20","reserved":0}`,
			},
			filename:    "syr20_E.db",
			minFileSize: 400000, // > 400kB
		},
		{
			name: "Download Watchtower",
			args: args{
				publJSON: `{"id":348729,"publicationRootKeyId":780,"mepsLanguageId":200,"publicationTypeId":14,"issueTagNumber":20210700,"title":"ਪਹਿਰਾਬੁਰਜ ਯਹੋਵਾਹ ਦੇ ਰਾਜ ਦੀ ਘੋਸ਼ਣਾ ਕਰਦਾ ਹੈ (ਸਟੱਡੀ)—2021","issueTitle":"ਪਹਿਰਾਬੁਰਜ, ਜੁਲਾਈ 2021","shortTitle":"ਪਹਿਰਾਬੁਰਜ (ਸਟੱਡੀ) (2021)","coverTitle":"ਅਧਿਐਨ ਲੇਖ: 30 ਅਗਸਤ–26 ਸਤੰਬਰ","undatedTitle":"ਪਹਿਰਾਬੁਰਜ - ਸਟੱਡੀ ਐਡੀਸ਼ਨ","undatedReferenceTitle":"ਪਹਿਰਾਬੁਰਜ (ਸਟੱਡੀ)","year":2021,"symbol":"w21","keySymbol":"w","reserved":0}`,
			},
			filename:    "w_PJ_202107.db",
			minFileSize: 600000, // > 400kB
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir("", "go-jwlm")
			assert.NoError(t, err)
			defer os.RemoveAll(tmp)

			dm := DownloadPublication(tt.args.publJSON, tmp)
			time.Sleep(250 * time.Millisecond)
			for progress := range dm.prgrsChan {
				assert.Greater(t, progress.BytesComplete, int64(1))
				if progress.Done {
					break
				}
			}
			if tt.wantErr {
				assert.Error(t, dm.err)
				return
			}

			info, err := os.Stat(filepath.Join(tmp, tt.filename))
			assert.NoError(t, err)
			assert.Greater(t, info.Size(), int64(tt.minFileSize))
		})
	}
}
