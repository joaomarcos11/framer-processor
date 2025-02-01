package service

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFile(t *testing.T) {
	dirOut := t.TempDir()
	fileName := "testfile.txt"

	file, err := CreateFile(fileName, dirOut)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer file.Close()

	expectedPath := filepath.Join(dirOut, fileName)
	if file.Name() != expectedPath {
		t.Errorf("expected file path %s, got %s", expectedPath, file.Name())
	}

	_, err = os.Stat(expectedPath)
	if os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", expectedPath)
	}
}

func TestZipFileByExtension(t *testing.T) {
	dir := t.TempDir()
	extension := ".txt"

	// Create test files
	files := []string{"file1.txt", "file2.txt", "file3.log"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	buf, err := ZipFileByExtension(dir, extension)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("expected non-empty buffer")
	}

	// Check if the zip contains the correct files
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("failed to read zip buffer: %v", err)
	}

	expectedFiles := map[string]bool{
		"file1.txt": false,
		"file2.txt": false,
	}

	for _, f := range r.File {
		if _, ok := expectedFiles[f.Name]; ok {
			expectedFiles[f.Name] = true
		} else {
			t.Errorf("unexpected file %s in zip", f.Name)
		}
	}

	for file, found := range expectedFiles {
		if !found {
			t.Errorf("expected file %s not found in zip", file)
		}
	}
}
