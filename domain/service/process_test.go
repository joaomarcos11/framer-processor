package service

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/mock"
)

// Mocking the usecases interfaces
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) DownloadFile(storage, file string) (*bytes.Buffer, error) {
	args := m.Called(storage, file)
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}

func (m *MockStorage) UploadFile(bucket, fileName string, file *bytes.Buffer) error {
	args := m.Called(bucket, fileName, file)
	return args.Error(0)
}

func (m *MockStorage) DeleteFile(bucket, fileName string) error {
	args := m.Called(bucket, fileName)
	return args.Error(0)
}

type MockMessager struct {
	mock.Mock
}

func (m *MockMessager) SendMessage(queue, message string) error {
	args := m.Called(queue, message)
	return args.Error(0)
}

func TestProcess(t *testing.T) {
	// mockStorage := new(MockStorage)
	// mockMessager := new(MockMessager)

	// Mocking the methods
	// mockMessager.On("SendMessage", "framer-status.fifo", "testfile.carregado").Return(nil)
	// mockStorage.On("DownloadFile", "teststorage", "testfile").Return(bytes.NewBuffer([]byte("file content")), nil)
	// mockMessager.On("SendMessage", "framer-status.fifo", "testfile.processando").Return(nil)
	// mockStorage.On("UploadFile", "fiap44-framer-images", "testfile.zip", mock.Anything).Return(nil)
	// mockMessager.On("SendMessage", "framer-status.fifo", "testfile.concluido").Return(nil)
	// mockStorage.On("DeleteFile", "fiap44-framer-videos", "testfile").Return(nil)

	// err := Process(mockStorage, mockMessager, "teststorage", "testfile")
	// assert.NoError(t, err)

	// // Test error cases
	// mockMessager.On("SendMessage", "framer-status.fifo", "testfile.carregado").Return(errors.New("send message error"))
	// err = Process(mockStorage, mockMessager, "teststorage", "testfile")
	// assert.Error(t, err)
	// assert.Equal(t, "failed to send message: send message error", err.Error())
}
