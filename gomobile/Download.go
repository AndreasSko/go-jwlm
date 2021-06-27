package gomobile

import (
	"context"

	"github.com/AndreasSko/go-jwlm/publication"
)

// DownloadManager keeps all the information of a running download, enabling it
// to check progress and also cancel the download if necessary
type DownloadManager struct {
	Progress  *DownloadProgress
	prgrsChan chan publication.Progress
	ctx       context.Context
	cancel    context.CancelFunc
	err       error
}

// DownloadProgress represents the progress of a running download
type DownloadProgress struct {
	Size           int64
	BytesComplete  int64
	BytesPerSecond float64
	Progress       float64
	Done           bool
	Canceled       bool
}

// CancelDownload cancels a running download
func (dm *DownloadManager) CancelDownload() {
	dm.cancel()
	dm.Progress.Canceled = true
}

// DownloadSuccessful indicates if the download has been successful
func (dm *DownloadManager) DownloadSuccessful() bool {
	return dm.Progress.Done && dm.err == nil && !dm.Progress.Canceled
}

// Error returns possible errors of a download as a string
func (dm *DownloadManager) Error() string {
	if dm.err != nil {
		return dm.err.Error()
	}
	return ""
}
