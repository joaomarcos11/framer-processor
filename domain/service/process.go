package service

import (
	"errors"
	"fmt"
	"framer-proc/domain/usecases"
)

func Process(strgSvc usecases.Storage, msgrSvc usecases.Messager, storage, file string) error {
	err := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "carregado"))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send message: %s", err))
	}

	obj, err := strgSvc.DownloadFile(storage, file)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to download file: %s", err))
	}
	defer obj.Close()

	videoFile, err := CreateFile(file, "/tmp")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create file: %s", err))
	}
	defer videoFile.Close()

	_, err = videoFile.ReadFrom(obj)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to read file: %s", err))
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "processando"))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send message: %s", err))
	}

	err = ExtractAndSaveFramesFromVideo(fmt.Sprintf("/tmp/%s", file), "/tmp")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get frames from video: %s", err))
	}

	buf, err := ZipFileByExtension("/tmp", ".jpg")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get generate zip file from video: %s", err))
	}

	err = strgSvc.UploadFile("fiap44-framer-images", fmt.Sprintf("%s.zip", file), &buf)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload file: %s", err))
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "concluido"))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send message: %s", err))
	}

	err = strgSvc.DeleteFile("fiap44-framer-videos", file)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete file: %s", err))
	}

	return nil
}
