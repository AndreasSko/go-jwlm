// +build !windows

package cmd

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/AndreasSko/go-jwlm/model"
	expect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func Test_merge(t *testing.T) {
	t.Parallel()

	tmp, err := ioutil.TempDir("", "go-jwlm")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)

	emptyFilename := filepath.Join(tmp, "empty.jwlibrary")
	leftFilename := filepath.Join(tmp, "left.jwlibrary")
	rightFilename := filepath.Join(tmp, "right.jwlibrary")
	mergedFilename := filepath.Join(tmp, "merged.jwlibrary")
	assert.NoError(t, emptyDB.ExportJWLBackup(emptyFilename))
	assert.NoError(t, leftDB.ExportJWLBackup(leftFilename))
	assert.NoError(t, rightDB.ExportJWLBackup(rightFilename))

	// Merge against empty DB and see if result is still the same
	RunCmdTest(t,
		func(t *testing.T, c *expect.Console) {
			_, err := c.ExpectString("🎉 Finished merging!")
			assert.NoError(t, err)
			_, err = c.ExpectEOF()
			assert.NoError(t, err)
		},
		func(t *testing.T, c *expect.Console) {
			merge(leftFilename, emptyFilename, mergedFilename,
				terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
			merged := &model.Database{}
			merged.ImportJWLBackup(mergedFilename)
			assert.True(t, leftDB.Equals(merged))
		})

	// Merge while selecting all right
	RunCmdTest(t,
		func(t *testing.T, c *expect.Console) {
			c.ExpectString("📑 Merging Bookmarks")
			c.SendLine(string(terminal.KeyArrowDown))

			c.ExpectString("🖍  Merging Markings")
			c.SendLine(string(terminal.KeyArrowDown))

			c.ExpectString("📝 Merging Notes")
			c.SendLine(string(terminal.KeyArrowDown))

			c.ExpectEOF()
		},
		func(t *testing.T, c *expect.Console) {
			merge(leftFilename, rightFilename, mergedFilename,
				terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
			merged := &model.Database{}
			merged.ImportJWLBackup(mergedFilename)
			assert.True(t, mergedAllRightDB.Equals(merged))
		})

	// Merge while selecting all left
	RunCmdTest(t,
		func(t *testing.T, c *expect.Console) {
			c.ExpectString("📑 Merging Bookmarks")
			c.SendLine("")

			c.ExpectString("🖍  Merging Markings")
			c.SendLine("")

			c.ExpectString("📝 Merging Notes")
			c.SendLine("")

			c.ExpectEOF()
		},
		func(t *testing.T, c *expect.Console) {
			merge(leftFilename, rightFilename, mergedFilename,
				terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
			merged := &model.Database{}
			merged.ImportJWLBackup(mergedFilename)
			assert.True(t, mergedAllLeftDB.Equals(merged))
		})
}

// https://github.com/AlecAivazis/survey/blob/master/survey_posix_test.go
func RunCmdTest(t *testing.T, procedure func(*testing.T, *expect.Console), test func(*testing.T, *expect.Console)) {
	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(t, c)
	}()

	test(t, c)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}

var emptyDB = &model.Database{}

var leftDB = &model.Database{
	BlockRange: []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{7, true},
			UserMarkID:   1,
		},
	},
	Bookmark: []*model.Bookmark{
		nil,
		{
			BookmarkID:            1,
			LocationID:            1,
			PublicationLocationID: 2,
			Title:                 "1. Mose 1:1",
			Snippet:               sql.NullString{"1 Am Anfang erschuf Gott Himmel und Erde.", true},
			BlockType:             2,
			BlockIdentifier:       sql.NullInt32{1, true},
		},
	},
	Location: []*model.Location{
		nil,
		{
			LocationID:    1,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{1, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 1", true},
		},
		{
			LocationID:   2,
			KeySymbol:    sql.NullString{"nwtsty", true},
			MepsLanguage: 2,
			LocationType: 1,
		},
	},
	Note: []*model.Note{
		nil,
		{
			NoteID:          1,
			GUID:            "92B192B4-B0B9-49B2-949F-7A8BA6406AF4",
			UserMarkID:      sql.NullInt32{1, true},
			LocationID:      sql.NullInt32{1, true},
			Title:           sql.NullString{"Am Anfang erschuf Gott Himmel und Erde.", true},
			Content:         sql.NullString{"📝 for left version", true},
			LastModified:    "2020-09-15T13:45:38+00:00",
			BlockType:       2,
			BlockIdentifier: sql.NullInt32{1, true},
		},
		{
			NoteID:       2,
			GUID:         "E36B34A0-B70F-4590-9D69-5887AB65A6D5",
			Title:        sql.NullString{"Same Note", true},
			Content:      sql.NullString{"This note is also available on the other side", true},
			LastModified: "2020-09-15T13:52:25+00:00",
			BlockType:    0,
		},
	},
	Tag: []*model.Tag{
		nil,
		{
			TagID:   1,
			TagType: 0,
			Name:    "Favorite",
		},
		{
			TagID:   2,
			TagType: 1,
			Name:    "Left",
		},
		{
			TagID:   3,
			TagType: 1,
			Name:    "Same",
		},
	},
	TagMap: []*model.TagMap{
		nil,
		{
			TagMapID: 1,
			NoteID:   sql.NullInt32{1, true},
			TagID:    2,
			Position: 0,
		},
		{
			TagMapID: 2,
			NoteID:   sql.NullInt32{2, true},
			TagID:    3,
			Position: 0,
		},
	},
	UserMark: []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   0,
			UserMarkGUID: "0D5523D9-F784-4B08-A86D-D517F4EB17DE",
			Version:      1,
		},
	},
}

var rightDB = &model.Database{
	BlockRange: []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{15, true},
			UserMarkID:   1,
		},
		{
			BlockRangeID: 2,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{7, true},
			UserMarkID:   2,
		},
		{
			BlockRangeID: 3,
			BlockType:    2,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{12, true},
			UserMarkID:   2,
		},
		{
			BlockRangeID: 4,
			BlockType:    2,
			Identifier:   2,
			StartToken:   sql.NullInt32{13, true},
			EndToken:     sql.NullInt32{26, true},
			UserMarkID:   3,
		},
	},
	Bookmark: []*model.Bookmark{
		nil,
		{
			BookmarkID:            1,
			LocationID:            1,
			PublicationLocationID: 2,
			Title:                 "1. Mose 2:1",
			Snippet:               sql.NullString{"2 So wurde die Erschaffung von Himmel und Erde und allem, was dazugehört, beendet. ", true},
			BlockType:             2,
			BlockIdentifier:       sql.NullInt32{1, true},
		},
	},
	Location: []*model.Location{
		nil,
		{
			LocationID:    1,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{2, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 2", true},
		},
		{
			LocationID:   2,
			KeySymbol:    sql.NullString{"nwtsty", true},
			MepsLanguage: 2,
			LocationType: 1,
		},
		{
			LocationID:    3,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{1, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 1", true},
		},
	},
	Note: []*model.Note{
		nil,
		{
			NoteID:          1,
			GUID:            "DE4A2CDA-9892-4A94-AF4B-22EBE05A05CA",
			UserMarkID:      sql.NullInt32{1, true},
			LocationID:      sql.NullInt32{1, true},
			Title:           sql.NullString{"So wurde die Erschaffung von Himmel und Erde und allem, was dazugehört, beendet.", true},
			Content:         sql.NullString{"📝 on the right side", true},
			LastModified:    "2020-09-15T13:47:56+00:00",
			BlockType:       2,
			BlockIdentifier: sql.NullInt32{1, true},
		},
		{
			NoteID:       2,
			GUID:         "E36B34A0-B70F-4590-9D69-5887AB65A6D5",
			Title:        sql.NullString{"Same Note", true},
			Content:      sql.NullString{"This note is also available on the other side. Though this one is newer 😏", true},
			LastModified: "2020-09-20T13:52:25+00:00",
			BlockType:    0,
		},
	},
	Tag: []*model.Tag{
		nil,
		{
			TagID:   1,
			TagType: 0,
			Name:    "Favorite",
		},
		{
			TagID:   2,
			TagType: 1,
			Name:    "Right",
		},
		{
			TagID:   3,
			TagType: 1,
			Name:    "Same",
		},
	},
	TagMap: []*model.TagMap{
		nil,
		{
			TagMapID: 1,
			NoteID:   sql.NullInt32{1, true},
			TagID:    2,
			Position: 0,
		},
		{
			TagMapID: 2,
			NoteID:   sql.NullInt32{2, true},
			TagID:    3,
			Position: 0,
		},
	},
	UserMark: []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   0,
			UserMarkGUID: "54C23825-AC1E-4890-8041-92B39C5E4B17",
			Version:      1,
		},
		{
			UserMarkID:   2,
			ColorIndex:   1,
			LocationID:   3,
			StyleIndex:   0,
			UserMarkGUID: "A098D2B0-7613-4676-9E96-442755105D9A",
		},
		{
			UserMarkID:   3,
			ColorIndex:   1,
			LocationID:   3,
			StyleIndex:   0,
			UserMarkGUID: "A86DECC8-69B1-4A43-A3A1-F1CA7B8E8388",
			Version:      1,
		},
	},
}

var mergedAllLeftDB = model.Database{
	BlockRange: []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{7, true},
			UserMarkID:   1,
		},
		{
			BlockRangeID: 2,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{15, true},
			UserMarkID:   2,
		},
		{
			BlockRangeID: 3,
			BlockType:    2,
			Identifier:   2,
			StartToken:   sql.NullInt32{13, true},
			EndToken:     sql.NullInt32{26, true},
			UserMarkID:   3,
		},
	},
	Bookmark: []*model.Bookmark{
		nil,
		{
			BookmarkID:            1,
			LocationID:            1,
			PublicationLocationID: 2,
			Title:                 "1. Mose 1:1",
			Snippet:               sql.NullString{"1 Am Anfang erschuf Gott Himmel und Erde.", true},
			BlockType:             2,
			BlockIdentifier:       sql.NullInt32{1, true},
		},
	},
	Location: []*model.Location{
		nil,
		{
			LocationID:    1,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{1, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 1", true},
		},
		{
			LocationID:   2,
			KeySymbol:    sql.NullString{"nwtsty", true},
			MepsLanguage: 2,
			LocationType: 1,
		},
		{
			LocationID:    3,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{2, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 2", true},
		},
	},
	Note: []*model.Note{
		nil,
		{
			NoteID:          1,
			GUID:            "92B192B4-B0B9-49B2-949F-7A8BA6406AF4",
			UserMarkID:      sql.NullInt32{1, true},
			LocationID:      sql.NullInt32{1, true},
			Title:           sql.NullString{"Am Anfang erschuf Gott Himmel und Erde.", true},
			Content:         sql.NullString{"📝 for left version", true},
			LastModified:    "2020-09-15T13:45:38+00:00",
			BlockType:       2,
			BlockIdentifier: sql.NullInt32{1, true},
		},
		{
			NoteID:       2,
			GUID:         "E36B34A0-B70F-4590-9D69-5887AB65A6D5",
			Title:        sql.NullString{"Same Note", true},
			Content:      sql.NullString{"This note is also available on the other side", true},
			LastModified: "2020-09-15T13:52:25+00:00",
			BlockType:    0,
		},
		{
			NoteID:          3,
			GUID:            "DE4A2CDA-9892-4A94-AF4B-22EBE05A05CA",
			UserMarkID:      sql.NullInt32{1, true},
			LocationID:      sql.NullInt32{1, true},
			Title:           sql.NullString{"So wurde die Erschaffung von Himmel und Erde und allem, was dazugehört, beendet.", true},
			Content:         sql.NullString{"📝 on the right side", true},
			LastModified:    "2020-09-15T13:47:56+00:00",
			BlockType:       2,
			BlockIdentifier: sql.NullInt32{1, true},
		},
	},
	Tag: []*model.Tag{
		nil,
		{
			TagID:   1,
			TagType: 0,
			Name:    "Favorite",
		},
		{
			TagID:   2,
			TagType: 1,
			Name:    "Left",
		},
		{
			TagID:   3,
			TagType: 1,
			Name:    "Same",
		},
		{
			TagID:   4,
			TagType: 1,
			Name:    "Right",
		},
	},
	TagMap: []*model.TagMap{
		nil,
		{
			TagMapID: 1,
			NoteID:   sql.NullInt32{1, true},
			TagID:    2,
			Position: 0,
		},
		{
			TagMapID: 2,
			NoteID:   sql.NullInt32{2, true},
			TagID:    3,
			Position: 0,
		},
		{
			TagMapID: 2,
			NoteID:   sql.NullInt32{3, true},
			TagID:    4,
			Position: 0,
		},
	},
	UserMark: []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   0,
			UserMarkGUID: "0D5523D9-F784-4B08-A86D-D517F4EB17DE",
			Version:      1,
		},
		{
			UserMarkID:   2,
			ColorIndex:   1,
			LocationID:   3,
			StyleIndex:   0,
			UserMarkGUID: "54C23825-AC1E-4890-8041-92B39C5E4B17",
			Version:      1,
		},
		{
			UserMarkID:   3,
			ColorIndex:   1,
			LocationID:   1,
			StyleIndex:   0,
			UserMarkGUID: "A86DECC-69B1-4A43-A3A1-F1CA7B8E8388",
			Version:      1,
		},
	},
}

var mergedAllRightDB = model.Database{
	BlockRange: []*model.BlockRange{
		nil,
		{
			BlockRangeID: 1,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{15, true},
			UserMarkID:   1,
		},
		{
			BlockRangeID: 2,
			BlockType:    2,
			Identifier:   1,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{7, true},
			UserMarkID:   2,
		},
		{
			BlockRangeID: 3,
			BlockType:    2,
			Identifier:   2,
			StartToken:   sql.NullInt32{0, true},
			EndToken:     sql.NullInt32{12, true},
			UserMarkID:   2,
		},
		{
			BlockRangeID: 4,
			BlockType:    2,
			Identifier:   2,
			StartToken:   sql.NullInt32{13, true},
			EndToken:     sql.NullInt32{26, true},
			UserMarkID:   3,
		},
	},
	Bookmark: []*model.Bookmark{
		nil,
		{
			BookmarkID:            1,
			LocationID:            3,
			PublicationLocationID: 2,
			Title:                 "1. Mose 2:1",
			Snippet:               sql.NullString{"2 So wurde die Erschaffung von Himmel und Erde und allem, was dazugehört, beendet. ", true},
			BlockType:             2,
			BlockIdentifier:       sql.NullInt32{1, true},
		},
	},
	Location: []*model.Location{
		nil,
		{
			LocationID:    1,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{1, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 1", true},
		},
		{
			LocationID:   2,
			KeySymbol:    sql.NullString{"nwtsty", true},
			MepsLanguage: 2,
			LocationType: 1,
		},
		{
			LocationID:    3,
			BookNumber:    sql.NullInt32{1, true},
			ChapterNumber: sql.NullInt32{2, true},
			KeySymbol:     sql.NullString{"nwtsty", true},
			MepsLanguage:  2,
			LocationType:  0,
			Title:         sql.NullString{"1. Mose 2", true},
		},
	},
	Note: []*model.Note{
		nil,
		{
			NoteID:          1,
			GUID:            "92B192B4-B0B9-49B2-949F-7A8BA6406AF4",
			UserMarkID:      sql.NullInt32{1, true},
			LocationID:      sql.NullInt32{1, true},
			Title:           sql.NullString{"Am Anfang erschuf Gott Himmel und Erde.", true},
			Content:         sql.NullString{"📝 for left version", true},
			LastModified:    "2020-09-15T13:45:38+00:00",
			BlockType:       2,
			BlockIdentifier: sql.NullInt32{1, true},
		},
		{
			NoteID:       2,
			GUID:         "E36B34A0-B70F-4590-9D69-5887AB65A6D5",
			Title:        sql.NullString{"Same Note", true},
			Content:      sql.NullString{"This note is also available on the other side. Though this one is newer 😏", true},
			LastModified: "2020-09-20T13:52:25+00:00",
			BlockType:    0,
		},
		{
			NoteID:          3,
			GUID:            "DE4A2CDA-9892-4A94-AF4B-22EBE05A05CA",
			UserMarkID:      sql.NullInt32{1, true},
			LocationID:      sql.NullInt32{1, true},
			Title:           sql.NullString{"So wurde die Erschaffung von Himmel und Erde und allem, was dazugehört, beendet.", true},
			Content:         sql.NullString{"📝 on the right side", true},
			LastModified:    "2020-09-15T13:47:56+00:00",
			BlockType:       2,
			BlockIdentifier: sql.NullInt32{1, true},
		},
	},
	Tag: []*model.Tag{
		nil,
		{
			TagID:   1,
			TagType: 0,
			Name:    "Favorite",
		},
		{
			TagID:   2,
			TagType: 1,
			Name:    "Left",
		},
		{
			TagID:   3,
			TagType: 1,
			Name:    "Same",
		},
		{
			TagID:   4,
			TagType: 1,
			Name:    "Right",
		},
	},
	TagMap: []*model.TagMap{
		nil,
		{
			TagMapID: 1,
			NoteID:   sql.NullInt32{1, true},
			TagID:    2,
			Position: 0,
		},
		{
			TagMapID: 2,
			NoteID:   sql.NullInt32{2, true},
			TagID:    3,
			Position: 0,
		},
		{
			TagMapID: 1,
			NoteID:   sql.NullInt32{3, true},
			TagID:    4,
			Position: 0,
		},
	},
	UserMark: []*model.UserMark{
		nil,
		{
			UserMarkID:   1,
			ColorIndex:   1,
			LocationID:   3, // 1. Mose 2
			StyleIndex:   0,
			UserMarkGUID: "54C23825-AC1E-4890-8041-92B39C5E4B17",
			Version:      1,
		},
		{
			UserMarkID:   2,
			ColorIndex:   1,
			LocationID:   1, // 1. Mose 1
			StyleIndex:   0,
			UserMarkGUID: "A098D2B0-7613-4676-9E96-442755105D9A",
		},
		{
			UserMarkID:   3,
			ColorIndex:   1,
			LocationID:   1, // 1. Mose 1
			StyleIndex:   0,
			UserMarkGUID: "A86DECC8-69B1-4A43-A3A1-F1CA7B8E8388",
			Version:      1,
		},
	},
}
