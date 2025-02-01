package service

import (
	"errors"
	"fmt"

	"github.com/filipeandrade6/framer-processor/domain/usecases"
)

func Process(strgSvc usecases.Storage, msgrSvc usecases.Messager, storage, file string) error {
	err := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "CARREGADO"))
	if err != nil {
		message := fmt.Sprintf("failed to send message: %s", err)
		return errors.New(message)
	}

	obj, err := strgSvc.DownloadFile(storage, file)
	if err != nil {
		message := fmt.Sprintf("failed to download file: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}
	defer obj.Close()

	videoFile, err := CreateFile(file, "/tmp")
	if err != nil {
		message := fmt.Sprintf("failed to create file: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}
	defer videoFile.Close()

	_, err = videoFile.ReadFrom(obj)
	if err != nil {
		message := fmt.Sprintf("failed to read file: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "PROCESSANDO"))
	if err != nil {
		message := fmt.Sprintf("failed to send message: %s", err)
		return errors.New(message)
	}

	err = ExtractAndSaveFramesFromVideo(fmt.Sprintf("/tmp/%s", file), "/tmp")
	if err != nil {
		message := fmt.Sprintf("failed to get frames from video: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}

	buf, err := ZipFileByExtension("/tmp", ".jpg")
	if err != nil {
		message := fmt.Sprintf("failed to get generate zip file from video: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}

	err = strgSvc.UploadFile("fiap44-framer-images", fmt.Sprintf("%s.zip", file), &buf)
	if err != nil {
		message := fmt.Sprintf("failed to upload file: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "CONCLUIDO"))
	if err != nil {
		message := fmt.Sprintf("failed to send message: %s", err)
		return errors.New(message)
	}

	err = strgSvc.DeleteFile("fiap44-framer-videos", file)
	if err != nil {
		message := fmt.Sprintf("failed to delete file: %s", err)
		msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, message))
		return errors.New(message)
	}

	return nil
}
