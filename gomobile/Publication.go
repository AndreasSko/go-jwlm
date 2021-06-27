package gomobile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AndreasSko/go-jwlm/model"
	"github.com/AndreasSko/go-jwlm/publication"
)

// PublicationLookup represents a lookup for a publication.
// It directly maps to publication.Lookup
type PublicationLookup struct {
	DocumentID     int
	KeySymbol      string
	IssueTagNumber int
	MepsLanguage   int
}

// LookupPublication looks up a publication from catalogDB located at dbPath
// and returns a JSON string representing the Publication
func LookupPublication(dbPath string, query *PublicationLookup) string {
	result, err := publication.LookupPublication(dbPath, publication.Lookup{
		DocumentID:     query.DocumentID,
		KeySymbol:      query.KeySymbol,
		IssueTagNumber: query.IssueTagNumber,
		MepsLanguage:   query.MepsLanguage,
	})
	if err != nil {
		return ""
	}

	jsn, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	return string(jsn)
}

// GetSnippet fetches the snippet related to a mergeConflict
func (mcw *MergeConflictsWrapper) GetSnippet(catalogDir, publDir string, conflictKey string) (string, error) {
	conflict, ok := mcw.conflicts[conflictKey]
	if !ok {
		return "", fmt.Errorf("conflict with key %s does not exist", conflictKey)
	}

	umbr, ok := conflict.Left.(*model.UserMarkBlockRange)
	if !ok {
		return "", fmt.Errorf("only UserMarkBlockRanges are supported for getting a snippet, %T given", conflict)
	}

	location := umbr.RelatedEntries(mcw.DBWrapper.left).Location

	publQuery := publication.Lookup{
		DocumentID:     int(location.DocumentID.Int32),
		KeySymbol:      location.KeySymbol.String,
		IssueTagNumber: location.IssueTagNumber,
		MepsLanguage:   location.MepsLanguage,
	}
	publ, err := publication.LookupPublication(catalogDir, publQuery)
	if err != nil {
		return "", fmt.Errorf("could not lookup publication: %w", err)
	}

	snippets, err := publication.GetSnippet(publDir, publ, *location, *umbr)
	if err != nil {
		return "", fmt.Errorf("could not get snippet: %w", err)
	}

	result, err := json.Marshal(snippets)
	if err != nil {
		return "", fmt.Errorf("could not marshal snippets to JSON: %w", err)
	}

	return string(result), nil
}

// GetPublicationPath generates the filename of the publication (given in JSON format)
// and checks if it exists in the publDir
func GetPublicationPath(publJSON string, publDir string) (string, error) {
	publ := publication.Publication{}
	if err := json.Unmarshal([]byte(publJSON), &publ); err != nil {
		return "", fmt.Errorf("could not unmarshal publication: %w", err)
	}

	return publication.GetPublicationPath(publ, publDir)
}

// DownloadPublication downloads the publication (given in JSON format) and saves it at dst. The
// returned DownloadManager allows to keep track and manage the running download
func DownloadPublication(publJSON string, dst string) *DownloadManager {
	publ := publication.Publication{}
	if err := json.Unmarshal([]byte(publJSON), &publ); err != nil {
		return &DownloadManager{
			Progress: &DownloadProgress{
				Done: true,
			},
			err: err,
		}
	}

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
			_, err := publication.DownloadPublication(ctx, dm.prgrsChan, publ, dst)
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
