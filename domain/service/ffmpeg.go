package service

import (
	"errors"
	"fmt"
	"os/exec"
)

func ExtractAndSaveFramesFromVideo(filePath, outDir string) error {
	_, err := exec.Command("/opt/ffmpeglib/ffmpeg", "-i", filePath, fmt.Sprintf("%s/frame_%%04d.jpg", outDir)).Output()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get video frames with ffmpeg: %v", err))
	}

	return nil
}
