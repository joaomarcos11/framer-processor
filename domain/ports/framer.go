package ports

type Framer interface {
	ExtractAndSaveFramesFromVideo(filePath, outDir string) error
}
