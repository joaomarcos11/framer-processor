package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for exec.Command
type MockCmd struct {
	mock.Mock
}

func (m *MockCmd) Output() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestExtractAndSaveFramesFromVideo(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		outDir    string
		mockError error
		wantError bool
	}{
		{
			name:      "success",
			filePath:  "test.mp4",
			outDir:    "output",
			mockError: nil,
			wantError: false,
		},
		{
			name:      "ffmpeg error",
			filePath:  "test.mp4",
			outDir:    "output",
			mockError: errors.New("ffmpeg error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmd := new(MockCmd)
			mockCmd.On("Output").Return([]byte{}, tt.mockError)

			// Replace exec.Command with our mock
			// execCommand := func(name string, arg ...string) *exec.Cmd {
			// 	return &exec.Cmd{}
			// }

			err := ExtractAndSaveFramesFromVideo(tt.filePath, tt.outDir)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockCmd.AssertExpectations(t)
		})
	}
}
