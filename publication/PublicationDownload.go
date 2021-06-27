package publication

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/codeclysm/extract/v3"
)

const publicationEndpoint = "https://app.jw-cdn.org/apis/pub-media/GETPUBMEDIALINKS?pub=%s&issue=%d&langwritten=%s&fileformat=jwpub"

// DownloadPublication downloads, unzips and stores the SQLiteDB of the publication in the directory at dst
func DownloadPublication(ctx context.Context, prgrs chan Progress, publ Publication, dst string) (string, error) {
	if prgrs != nil {
		defer close(prgrs)
	}

	// Create tmp folder and place all files there
	tmp, err := ioutil.TempDir("", "go-jwlm")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmp)

	url, err := getPublicationURL(ctx, publ)
	if err != nil {
		return "", fmt.Errorf("could not get publicationURL: %w", err)
	}

	filename, err := download(ctx, prgrs, url, tmp)
	if err != nil {
		return "", fmt.Errorf("could not download publication from %s: %w", url, err)
	}

	// Extract JWPub file (first layer)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %w", filename, err)
	}
	buffer := bytes.NewBuffer(data)
	err = extract.Zip(ctx, buffer, tmp, nil)
	if err != nil {
		return "", fmt.Errorf("could not extract first layer of publication: %w", err)
	}

	// Extract contents (second layer)
	filename = filepath.Join(tmp, "contents")
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %w", filename, err)
	}
	buffer = bytes.NewBuffer(data)
	err = extract.Zip(ctx, buffer, tmp, nil)
	if err != nil {
		return "", fmt.Errorf("could not extract first layer of publication: %w", err)
	}

	// List all files in tmp and pick *.db file to move to dst (the rest we don't need)
	files, err := ioutil.ReadDir(tmp)
	if err != nil {
		return "", fmt.Errorf("could not list files in tmp: %w", err)
	}
	filename = ""
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".db" {
			filename = file.Name()
		}
	}
	if filename == "" {
		return "", fmt.Errorf("could not find publication .db file")
	}

	dst = filepath.Join(dst, filename)
	os.Rename(filepath.Join(tmp, filename), dst)

	return dst, nil
}

// getPublicationURL looks up the URL for downloading the given publication
func getPublicationURL(ctx context.Context, publ Publication) (string, error) {
	language, err := lookupMepsLanguage(publ.MepsLanguageID)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(publicationEndpoint, publ.KeySymbol.String, publ.IssueTagNumber, language.Symbol)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error while creating new request for %s: %w", url, err)
	}
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not download publication manifest from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response body for publication manifest %s: %w", url, err)
	}

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return "", fmt.Errorf("could not parse publication manifest from %s: %w", url, err)
	}

	var publicationURL string
	publicationURL, ok := jsonParsed.Search("files", language.Symbol, "JWPUB", "file", "url").Index(0).Data().(string)
	if !ok {
		return "", fmt.Errorf("could not get url from publication manifest %s from %s", body, url)
	}

	return publicationURL, nil
}
