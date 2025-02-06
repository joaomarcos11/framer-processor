package ports

import (
	"bytes"
	"io"
)

type Filer interface {
	CreateFileWithContents(fileName, dirOut string, contents io.ReadCloser) error
	ZipFileByExtension(dir, extension string) (bytes.Buffer, error)
}

// type File interface {
// 	io.ReadCloser
// 	io.Reader
// }
