package service

import (
	"fmt"
	"os/exec"
)

func ExtractAndSaveFramesFromVideo(filePath, outDir string) error {
	_, err := exec.Command("/opt/ffmpeglib/ffmpeg", "-i", filePath, fmt.Sprintf("%s/frame_%%04d.jpg", outDir)).Output()
	if err != nil {
		return fmt.Errorf("extract frames from video: %w", err)
	}

	return nil
}
