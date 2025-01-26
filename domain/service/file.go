package service

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateFile(fileName, dirOut string) (*os.File, error) {
	file, err := os.Create(fmt.Sprintf("%s/%s", dirOut, fileName))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to crate file %s in directory %s: %s", fileName, dirOut, err))
	}

	return file, nil
}

func ZipFileByExtension(dir, extension string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	zipFile := io.Writer(&buf)

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	filesInDir, err := os.ReadDir(dir)
	if err != nil {
		return bytes.Buffer{}, errors.New(fmt.Sprintf("failed to read directory: %s", err))
	}

	files := []string{}
	for _, file := range filesInDir {
		if filepath.Ext(file.Name()) == extension {
			files = append(files, file.Name())
		}
	}

	if len(files) == 0 {
		return bytes.Buffer{}, errors.New("no frames extracted from video")
	}

	for _, file := range files {
		fileToZip, err := os.Open(fmt.Sprintf("/tmp/%s", file))
		if err != nil {
			return bytes.Buffer{}, errors.New(fmt.Sprintf("failed to open file %s: %v", file, err))
		}
		defer fileToZip.Close()

		zipEntry, err := zipWriter.Create(file)
		if err != nil {
			return bytes.Buffer{}, errors.New(fmt.Sprintf("failed to add file entry from file %s to the zip archive: %v", file, err))
		}

		_, err = io.Copy(zipEntry, fileToZip)
		if err != nil {
			return bytes.Buffer{}, errors.New(fmt.Sprintf("failed to write the file contents from file %s to the zip archive: %v", file, err))
		}
	}

	return buf, nil
}
