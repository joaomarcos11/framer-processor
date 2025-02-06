package errors

import "errors"

var (
	ErrEmptyFile     = errors.New("empty file")
	ErrFileTooBig    = errors.New("file too big")
	ErrReadDir       = errors.New("failed to read directory")
	ErrSendMessage   = errors.New("failed to send message")
	ErrCreateFile    = errors.New("failed to create file with contents")
	ErrExtractFrames = errors.New("failed to extract frames from video")
	ErrOpeningFile   = errors.New("failed to open file")
	ErrZipFile       = errors.New("failed to create zip file")
	ErrDownloadFile  = errors.New("failed to download file")
	ErrReadFile      = errors.New("failed to read file")
	ErrUploadFile    = errors.New("failed to upload file")
	ErrDeleteFile    = errors.New("failed to delete file")
)
