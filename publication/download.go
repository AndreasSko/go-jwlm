package publication

import (
	"context"
	"fmt"
	"time"

	"github.com/cavaliercoder/grab"
)

// Progress represents the progress of a running download
type Progress struct {
	Size           int64
	BytesComplete  int64
	BytesPerSecond float64
	Progress       float64
	Duration       time.Duration
	ETA            time.Time
	Done           bool
}

// download downloads a file from url and stores it at dst.
// The prgrs channel informs about the progress of the download.
func download(ctx context.Context, prgrs chan Progress, url string, dst string) (string, error) {
	client := grab.NewClient()
	req, err := grab.NewRequest(dst, url)
	if err != nil {
		return "", fmt.Errorf("could not create request for %s: %w", url, err)
	}
	req = req.WithContext(ctx)
	resp := client.Do(req)

	progress := Progress{}

	// Send a progress over the prgrsChan every 250 milliseconds
	t := time.NewTicker(250 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		progress := Progress{
			Size:           resp.Size(),
			BytesComplete:  resp.BytesComplete(),
			BytesPerSecond: resp.BytesPerSecond(),
			Progress:       resp.Progress(),
			Duration:       resp.Duration(),
			ETA:            resp.ETA(),
		}
		select {
		case <-t.C:
			select {
			case prgrs <- progress:
				continue
			default:
				continue
			}
		case <-resp.Done:
			break Loop
		}
	}
	if err := resp.Err(); err != nil {
		progress.Done = true
		select {
		case prgrs <- progress:
		default:
		}
		return "", fmt.Errorf("could not download %s: %w", url, err)
	}

	return resp.Filename, nil
}
