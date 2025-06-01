package publication

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCatalogNeedsUpdate(t *testing.T) {
	tmp := t.TempDir()

	assert.True(t, CatalogNeedsUpdate("not-valid-path"))

	filePath := filepath.Join(tmp, "catalog.db")
	_, err := os.Create(filePath)
	assert.NoError(t, err)

	assert.False(t, CatalogNeedsUpdate(filePath))

	os.Chtimes(filePath, time.Now(), time.Now().Add(-time.Hour*24*31))
	assert.True(t, CatalogNeedsUpdate(filePath))
}

func TestCatalogExists(t *testing.T) {
	tmp := t.TempDir()

	filePath := filepath.Join(tmp, "catalog.db")
	_, err := os.Create(filePath)
	assert.NoError(t, err)

	assert.False(t, CatalogExists("not-valid-path"))
	assert.True(t, CatalogExists(filePath))
}

func TestCatalogSize(t *testing.T) {
	assert.Equal(t, int64(77824), CatalogSize(filepath.Join("testdata", "catalog.db")))
	assert.Equal(t, int64(0), CatalogSize("not-valid-path"))
}

func Test_DownloadCatalogRealLife(t *testing.T) {
	tmp := t.TempDir()

	prgrs := make(chan Progress)
	done := make(chan struct{})
	go func() {
		assert.NoError(t, DownloadCatalog(context.Background(), prgrs, filepath.Join(tmp, "catalog.db")))
		done <- struct{}{}
	}()
	for progress := range prgrs {
		assert.IsType(t, Progress{}, progress)
		assert.NotEqual(t, Progress{}, progress)
	}
	<-done

	info, err := os.Stat(filepath.Join(tmp, "catalog.db"))
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100000000)) // catalog should by > 100MB
}

func Test_DownloadCatalog(t *testing.T) {
	tmp := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.String(), "manifest.json") {
			rw.Write([]byte(`{"version": 1, "current": "164a1c4b-4dbd-4909-8f88-8e7a18c562f2"}`))
		} else {
			data, err := os.ReadFile(filepath.Join("testdata", "catalog.db.gz"))
			assert.NoError(t, err)
			rw.Write(data)
		}
	}))
	defer server.Close()

	ManifestURL = server.URL + "/catalogs/publications/v4/manifest.json"
	CatalogURL = server.URL + "/catalogs/publications/v4/%s/catalog.db.gz"

	// Download without progress channel
	err := DownloadCatalog(context.TODO(), nil, filepath.Join(tmp, "catalog.db"))
	assert.NoError(t, err)

	assert.Equal(t, "7ebe98db8b5edd1ab901b7d6b43647fd35790b2a332c43739efdf9383d590651",
		hashFile(filepath.Join(tmp, "catalog.db")))

	// Old catalog should be overridden & test with progress channel
	prgrs := make(chan Progress)
	done := make(chan struct{})
	go func() {
		err := DownloadCatalog(context.Background(), prgrs, filepath.Join(tmp, "catalog.db"))
		assert.NoError(t, err)
		done <- struct{}{}
	}()
	for progress := range prgrs {
		assert.IsType(t, Progress{}, progress)
		assert.NotEqual(t, Progress{}, progress)
	}
	<-done

	// Check if timeout works
	ctx, cancle := context.WithTimeout(context.Background(), time.Duration(1))
	defer cancle()
	err = DownloadCatalog(ctx, nil, filepath.Join(tmp, "catalog.db"))
	assert.Error(t, err)
}

func Test_fetchManifest(t *testing.T) {
	var tests = []struct {
		input       string
		expected    catalogManifest
		expectError bool
	}{
		{
			input: `{"version": 1, "current": "164a1c4b-4dbd-4909-8f88-8e7a18c562f2"}`,
			expected: catalogManifest{
				Version: 1,
				Current: "164a1c4b-4dbd-4909-8f88-8e7a18c562f2",
			},
		},
		{
			input:       "ERROR",
			expected:    catalogManifest{},
			expectError: true,
		},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(test.input))
		}))
		defer server.Close()

		ManifestURL = server.URL
		res, err := fetchManifest(context.Background())
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expected, res)
	}
}

func hashFile(path string) string {
	f, _ := os.Open(path)
	hasher := sha256.New()
	io.Copy(hasher, f)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
