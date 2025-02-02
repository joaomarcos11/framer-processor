package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateFile(fileName, dirOut string) (*os.File, error) {
	file, err := os.Create(fmt.Sprintf("%s/%s", dirOut, fileName))
	if err != nil {
		err = fmt.Errorf("failed to create file %s in directory %s: %w", fileName, dirOut, err)
		return nil, err
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
		err = fmt.Errorf("failed to read directory: %w", err)
		return bytes.Buffer{}, err
	}

	files := []string{}
	for _, file := range filesInDir {
		if filepath.Ext(file.Name()) == extension {
			files = append(files, file.Name())
		}
	}

	if len(files) == 0 {
		err = fmt.Errorf("no frames extracted from video: %w", err)
		return bytes.Buffer{}, err
	}

	for _, file := range files {
		fileToZip, err := os.Open(fmt.Sprintf("/tmp/%s", file))
		if err != nil {
			err = fmt.Errorf("failed to open file %s: %w", file, err)
			return bytes.Buffer{}, err
		}
		defer fileToZip.Close()

		zipEntry, err := zipWriter.Create(file)
		if err != nil {
			err = fmt.Errorf("failed to add file entry from file %s to the zip archive: %w", file, err)
			return bytes.Buffer{}, err
		}

		_, err = io.Copy(zipEntry, fileToZip)
		if err != nil {
			err = fmt.Errorf("failed to write the file contents from file %s to the zip archive: %w", file, err)
			return bytes.Buffer{}, err
		}
	}

	return buf, nil
}
