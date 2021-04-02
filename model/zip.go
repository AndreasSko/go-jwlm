package model

import (
	"archive/zip"
	"bytes"
	"io"
	"time"

	"github.com/pkg/errors"
)

// https://golangcode.com/create-zip-files-in-go/
func zipFiles(filename string, files []string) error {
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	// Add files to zip
	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			zipWriter.Close()
			return err
		}
	}

	zipWriter.Close()
	err := GetPersistence().WriteFile(filename, zipBuffer.Bytes())
	if err != nil {

		return errors.Wrap(err, "Error storing zip")
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	name, data, err := GetPersistence().GetFile(filename)
	if err != nil {
		return err
	}

	header := zip.FileHeader{}
	header.Name = name
	header.Method = zip.Deflate
	header.Modified = time.Now()

	writer, err := zipWriter.CreateHeader(&header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(data))
	return err
}
