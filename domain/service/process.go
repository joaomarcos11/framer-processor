package service

import (
	"fmt"
	"log/slog"

	"github.com/filipeandrade6/framer-processor/domain/errors"
	"github.com/filipeandrade6/framer-processor/domain/usecases"
)

func Process(strgSvc usecases.Storage, msgrSvc usecases.Messager, storage, file string) error {
	err := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "CARREGADO"))
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err)
		slog.Error("send message", "err", err)
		return err
	}

	obj, err := strgSvc.DownloadFile(storage, file)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrDownloadFile, err)
		slog.Error("download file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}
	defer obj.Close()

	videoFile, err := CreateFile(file, "/tmp")
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrCreateFile, err)
		slog.Error("create file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}
	defer videoFile.Close()

	_, err = videoFile.ReadFrom(obj)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrReadFile, err)
		slog.Error("read file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "PROCESSANDO"))
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err)
		slog.Error("send message", "err", err)
		return err
	}

	err = ExtractAndSaveFramesFromVideo(fmt.Sprintf("/tmp/%s", file), "/tmp")
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrExtractFrames, err)
		slog.Error("extract frames", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}

	buf, err := ZipFileByExtension("/tmp", ".jpg")
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrZipFile, err)
		slog.Error("zip files", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}

	err = strgSvc.UploadFile("fiap44-framer-images", fmt.Sprintf("%s.zip", file), &buf)
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrUploadFile, err)
		slog.Error("upload file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "CONCLUIDO"))
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrSendMessage, err)
		slog.Error("send message", "err", err)
		return err
	}

	err = strgSvc.DeleteFile("fiap44-framer-videos", file)
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrDeleteFile, err)
		slog.Error("delete file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err)
		}
		return err
	}

	return nil
}
