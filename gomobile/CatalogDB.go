package gomobile

import (
	"context"

	"github.com/AndreasSko/go-jwlm/publication"
)

// DownloadCatalog downloads the newest catalog.db and saves it at dst. The
// returned DownloadManager allows to keep track and manage the running download
func DownloadCatalog(dst string) *DownloadManager {
	ctx, cancel := context.WithCancel(context.Background())
	dm := &DownloadManager{
		Progress:  &DownloadProgress{},
		prgrsChan: make(chan publication.Progress),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Start download in sub-goroutine, while monitoring its progress
	go func() {
		done := make(chan struct{})
		go func() {
			err := publication.DownloadCatalog(dm.ctx, dm.prgrsChan, dst)
			if err != nil {
				dm.err = err
			}
			done <- struct{}{}
		}()
		for progress := range dm.prgrsChan {
			dm.Progress.Size = progress.Size
			dm.Progress.BytesComplete = progress.BytesComplete
			dm.Progress.BytesPerSecond = progress.BytesPerSecond
			dm.Progress.Progress = progress.Progress
		}
		<-done
		dm.Progress.Done = true
	}()

	return dm
}

// CatalogNeedsUpdate checks if catalog.db located at path is still up-to-date.
// For now it just makes sure that it is younger than one month.
// If it can't find a file at path, it returns true
func CatalogNeedsUpdate(path string) bool {
	return publication.CatalogNeedsUpdate(path)
}

// CatalogExists checks if catalog.db exists at path
func CatalogExists(path string) bool {
	return publication.CatalogExists(path)
}

// CatalogSize returns the size of the catalog.db at path
func CatalogSize(path string) int64 {
	return publication.CatalogSize(path)
}
