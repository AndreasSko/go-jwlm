// +build !windows

package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/AndreasSko/go-jwlm/storage"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func Test_compare(t *testing.T) {
	t.Parallel()

	tmp, err := ioutil.TempDir("", "go-jwlm")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)

	emptyFilename := filepath.Join(tmp, "empty.jwlibrary")
	leftFilename := filepath.Join(tmp, "left.jwlibrary")
	rightFilename := filepath.Join(tmp, "right.jwlibrary")
	assert.NoError(t, storage.ExportJWLBackup(emptyDB, emptyFilename))
	assert.NoError(t, storage.ExportJWLBackup(leftDB, leftFilename))
	assert.NoError(t, storage.ExportJWLBackup(rightDB, rightFilename))

	RunCmdTest(t,
		func(t *testing.T, c *expect.Console) {
			_, err := c.ExpectString("Backups are NOT equal")
			assert.NoError(t, err)
			_, err = c.ExpectEOF()
			assert.NoError(t, err)
		},
		func(t *testing.T, c *expect.Console) {
			compare(leftFilename, emptyFilename, terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
			time.Sleep(time.Millisecond * 150) // So it does not finish before go-expect finished
		})

	RunCmdTest(t,
		func(t *testing.T, c *expect.Console) {
			_, err := c.ExpectString("✅ Backups are equal")
			assert.NoError(t, err)
			_, err = c.ExpectEOF()
			assert.NoError(t, err)
		},
		func(t *testing.T, c *expect.Console) {
			compare(leftFilename, leftFilename, terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
			time.Sleep(time.Millisecond * 150) // So it does not finish before go-expect finished
		})

	RunCmdTest(t,
		func(t *testing.T, c *expect.Console) {
			_, err := c.ExpectString("Backups are equal")
			assert.NoError(t, err)
			_, err = c.ExpectEOF()
			assert.NoError(t, err)
		},
		func(t *testing.T, c *expect.Console) {
			compare(rightFilename, rightFilename, terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
			time.Sleep(time.Millisecond * 150) // So it does not finish before go-expect finished
		})
}
