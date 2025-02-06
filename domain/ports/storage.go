package ports

import "io"

type Storager interface {
	DownloadFile(repository, fileName string) (io.ReadCloser, error)
	UploadFile(repository, fileName string, data io.Reader) error
	DeleteFile(repository, fileName string) error
}
