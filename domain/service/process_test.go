package service

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations of the ports
type MockStorager struct {
	mock.Mock
}

func (m *MockStorager) DownloadFile(storage, file string) (io.ReadCloser, error) {
	args := m.Called(storage, file)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// func (m *MockStorager) UploadFile(bucket, file string, buf *bytes.Buffer) error {
func (m *MockStorager) UploadFile(bucket, file string, buf io.Reader) error {
	args := m.Called(bucket, file, buf)
	return args.Error(0)
}

func (m *MockStorager) DeleteFile(bucket, file string) error {
	args := m.Called(bucket, file)
	return args.Error(0)
}

type MockMessager struct {
	mock.Mock
}

func (m *MockMessager) SendMessage(queue, message string) error {
	args := m.Called(queue, message)
	return args.Error(0)
}

type MockFiler struct {
	mock.Mock
}

func (m *MockFiler) CreateFileWithContents(file, path string, contents io.ReadCloser) error {
	args := m.Called(file, path, contents)
	return args.Error(0)
}

func (m *MockFiler) ZipFileByExtension(path, ext string) (bytes.Buffer, error) {
	args := m.Called(path, ext)
	return args.Get(0).(bytes.Buffer), args.Error(1)
}

type MockFramer struct {
	mock.Mock
}

func (m *MockFramer) ExtractAndSaveFramesFromVideo(file, path string) error {
	args := m.Called(file, path)
	return args.Error(0)
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name          string
		size          int64
		setupMocks    func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer)
		expectedError error
	}{
		{
			name: "empty file",
			size: 0,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.empty file").Return(nil)
			},
			expectedError: errors.New("empty file"),
		},
		{
			name: "empty file and failed to send message",
			size: 0,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.empty file").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("empty file"),
		},
		{
			name: "file too big",
			size: 2e7 + 1,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.file too big").Return(nil)
			},
			expectedError: errors.New("file too big"),
		},
		{
			name: "file too big and failed to send message",
			size: 2e7 + 1,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.file too big").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("file too big"),
		},
		{
			name: "failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to send message"),
		},
		{
			name: "failed to download file",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(new(os.File), errors.New("failed to download file"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to download file: failed to download file").Return(nil)
			},
			expectedError: errors.New("failed to download file"),
		},
		{
			name: "failed to download file and failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(new(os.File), errors.New("failed to download file"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to download file: failed to download file").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to download file"),
		},
		{
			name: "filed to create file with contents",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(errors.New("failed to create file with contents"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to create file with contents: failed to create file with contents").Return(nil)
			},
			expectedError: errors.New("failed to create file with contents"),
		},
		{
			name: "filed to create file with contents and failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(errors.New("failed to create file with contents"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to create file with contents: failed to create file with contents").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to create file with contents"),
		},
		{
			name: "failed to send message after download file",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to send message"),
		},
		{
			name: "failed to extract and save frames from video",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(errors.New("failed to extract and save frames from video"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to extract frames from video:failed to extract and save frames from video").Return(nil)
			},
			expectedError: errors.New("failed to extract and save frames from video"),
		},
		{
			name: "failed to extract and save frames from video and failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(errors.New("failed to extract and save frames from video"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to extract frames from video:failed to extract and save frames from video").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to extract and save frames from video"),
		},
		{
			name: "failed to zip file by extension",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, errors.New("failed to zip file by extension"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to create zip file:failed to zip file by extension").Return(nil)
			},
			expectedError: errors.New("failed to zip file by extension"),
		},
		{
			name: "failed to zip file by extension and failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, errors.New("failed to zip file by extension"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to create zip file:failed to zip file by extension").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to zip file by extension"),
		},
		{
			name: "failed to upload file",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, nil)
				strgSvc.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(errors.New("failed to upload file"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to upload file:failed to upload file").Return(nil)
			},
			expectedError: errors.New("failed to upload file"),
		},
		{
			name: "failed to upload file and failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, nil)
				strgSvc.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(errors.New("failed to upload file"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to upload file:failed to upload file").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to upload file"),
		},
		{
			name: "failed to send message after upload file",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, nil)
				strgSvc.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CONCLUIDO").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to send message"),
		},
		{
			name: "failed to delete file after process",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, nil)
				strgSvc.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CONCLUIDO").Return(nil)
				strgSvc.On("DeleteFile", "fiap44-framer-videos", "testfile").Return(errors.New("failed to delete file"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to delete file:failed to delete file").Return(nil)
			},
			expectedError: errors.New("failed to delete file"),
		},
		{
			name: "failed to delete file after process and failed to send message",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, nil)
				strgSvc.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CONCLUIDO").Return(nil)
				strgSvc.On("DeleteFile", "fiap44-framer-videos", "testfile").Return(errors.New("failed to delete file"))
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.FALHA.failed to delete file:failed to delete file").Return(errors.New("failed to send message"))
			},
			expectedError: errors.New("failed to delete file"),
		},
		{
			name: "successful process",
			size: 1e6,
			setupMocks: func(strgSvc *MockStorager, msgrSvc *MockMessager, filer *MockFiler, framer *MockFramer) {
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CARREGADO").Return(nil)
				strgSvc.On("DownloadFile", "teststorage", "testfile").Return(io.NopCloser(bytes.NewReader([]byte("file content"))), nil)
				filer.On("CreateFileWithContents", "testfile", "/tmp", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.PROCESSANDO").Return(nil)
				framer.On("ExtractAndSaveFramesFromVideo", "/tmp/testfile", "/tmp").Return(nil)
				filer.On("ZipFileByExtension", "/tmp", ".jpg").Return(bytes.Buffer{}, nil)
				strgSvc.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(nil)
				msgrSvc.On("SendMessage", "framer-status.fifo", "testfile.CONCLUIDO").Return(nil)
				strgSvc.On("DeleteFile", "fiap44-framer-videos", "testfile").Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strgSvc := new(MockStorager)
			msgrSvc := new(MockMessager)
			filer := new(MockFiler)
			framer := new(MockFramer)

			tt.setupMocks(strgSvc, msgrSvc, filer, framer)

			err := Process(strgSvc, msgrSvc, filer, framer, "teststorage", "testfile", tt.size)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			strgSvc.AssertExpectations(t)
			msgrSvc.AssertExpectations(t)
			filer.AssertExpectations(t)
			framer.AssertExpectations(t)
		})
	}
}
