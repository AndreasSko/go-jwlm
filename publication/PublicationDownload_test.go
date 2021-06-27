package publication

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadPublication(t *testing.T) {
	type args struct {
		publ Publication
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
				publ: Publication{
					ID:                    317733,
					PublicationRootKeyID:  786,
					MepsLanguageID:        0,
					PublicationTypeID:     2,
					IssueTagNumber:        0,
					Title:                 "2020 Service Year Report of Jehovah’s Witnesses Worldwide",
					ShortTitle:            "2020 Service Year Report",
					UndatedReferenceTitle: sql.NullString{"2020 Service Year Report", true},
					Year:                  2020,
					Symbol:                "syr20",
					KeySymbol:             sql.NullString{"syr20", true},
				},
			},
			filename:    "syr20_E.db",
			minFileSize: 400000, // > 400kB
		},
		{
			name: "Download Watchtower",
			args: args{
				publ: Publication{
					ID:                    348729,
					PublicationRootKeyID:  780,
					MepsLanguageID:        200,
					PublicationTypeID:     14,
					IssueTagNumber:        20210700,
					Title:                 "ਪਹਿਰਾਬੁਰਜ ਯਹੋਵਾਹ ਦੇ ਰਾਜ ਦੀ ਘੋਸ਼ਣਾ ਕਰਦਾ ਹੈ (ਸਟੱਡੀ)—2021",
					IssueTitle:            sql.NullString{"ਪਹਿਰਾਬੁਰਜ, ਜੁਲਾਈ 2021", true},
					ShortTitle:            "ਪਹਿਰਾਬੁਰਜ (ਸਟੱਡੀ) (2021)",
					CoverTitle:            sql.NullString{"ਅਧਿਐਨ ਲੇਖ: 30 ਅਗਸਤ–26 ਸਤੰਬਰ", true},
					UndatedTitle:          sql.NullString{"ਪਹਿਰਾਬੁਰਜ - ਸਟੱਡੀ ਐਡੀਸ਼ਨ", true},
					UndatedReferenceTitle: sql.NullString{"ਪਹਿਰਾਬੁਰਜ (ਸਟੱਡੀ)", true},
					Year:                  2021,
					Symbol:                "w21",
					KeySymbol:             sql.NullString{"w", true},
				},
			},
			filename:    "w_PJ_202107.db",
			minFileSize: 600000, // > 600kB
		},
		{
			name: "Wrong download",
			args: args{
				publ: Publication{
					ID:                   12345678910,
					PublicationRootKeyID: 12345678910,
					MepsLanguageID:       12345678910,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir("", "go-jwlm")
			assert.NoError(t, err)
			defer os.RemoveAll(tmp)

			prgrs := make(chan Progress)
			done := make(chan struct{})
			var path string
			go func() {
				path, err = DownloadPublication(context.Background(), prgrs, tt.args.publ, tmp)
				if tt.wantErr {
					assert.Error(t, err)
					done <- struct{}{}
					return
				}
				assert.NoError(t, err)
				done <- struct{}{}
			}()
			for progress := range prgrs {
				assert.IsType(t, Progress{}, progress)
				assert.NotEqual(t, Progress{}, progress)
			}
			<-done

			info, err := os.Stat(path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Greater(t, info.Size(), int64(tt.minFileSize))

			filename := filepath.Base(path)
			assert.Equal(t, tt.filename, filename)
		})
	}
}

func Test_getPublicationURL(t *testing.T) {
	type args struct {
		ctx  context.Context
		publ Publication
	}
	tests := []struct {
		name     string
		args     args
		want     string
		errorMsg string
	}{
		{
			name: "Get `Draw Close to Jehovah`",
			args: args{
				ctx: context.TODO(),
				publ: Publication{
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
				},
			},
			want: "https://download-a.akamaihd.net/files/media_publication/36/cl_E.jwpub",
		},
		{
			name: "Get Watchtower issue in Chinese",
			args: args{
				ctx: context.TODO(),
				publ: Publication{
					ID:                    305097,
					PublicationRootKeyID:  780,
					MepsLanguageID:        43,
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
				},
			},
			want: "https://download-a.akamaihd.net/files/media_periodical/92/w_CH_202102.jwpub",
		},
		{
			name: "Get with invalid MepsLanguageID",
			args: args{
				ctx: context.TODO(),
				publ: Publication{
					ID:                    67,
					PublicationRootKeyID:  64,
					MepsLanguageID:        12345,
					PublicationTypeID:     2,
					IssueTagNumber:        0,
					Title:                 "Draw Close to Jehovah",
					ShortTitle:            "Close to Jehovah",
					UndatedReferenceTitle: sql.NullString{"Close to Jehovah", true},
					Year:                  2014,
					Symbol:                "cl",
					KeySymbol:             sql.NullString{"cl", true},
				},
			},
			errorMsg: "could not find language symbol for mepsLanguageID 12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPublicationURL(tt.args.ctx, tt.args.publ)
			if tt.errorMsg != "" {
				assert.Error(t, err)
				assert.EqualErrorf(t, err, tt.errorMsg, "")
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
