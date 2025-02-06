package service

import (
	"fmt"
	"log/slog"

	"github.com/filipeandrade6/framer-processor/domain/errors"
	"github.com/filipeandrade6/framer-processor/domain/ports"
)

func Process(strgSvc ports.Storager, msgrSvc ports.Messager, filer ports.Filer, framer ports.Framer, storage, file string, size int64) error {
	if size == 0 {
		err := errors.ErrEmptyFile
		slog.Error("download file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

	if size > 2e7 {
		err := errors.ErrFileTooBig
		slog.Error("download file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

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
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

	err = filer.CreateFileWithContents(file, "/tmp", obj)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrCreateFile, err)
		slog.Error("create file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

	err = msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.%s", file, "PROCESSANDO"))
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.ErrSendMessage, err)
		slog.Error("send message", "err", err)
		return err
	}

	err = framer.ExtractAndSaveFramesFromVideo(fmt.Sprintf("/tmp/%s", file), "/tmp")
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrExtractFrames, err)
		slog.Error("extract frames", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

	buf, err := filer.ZipFileByExtension("/tmp", ".jpg")
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrZipFile, err)
		slog.Error("zip files", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

	err = strgSvc.UploadFile("fiap44-framer-images", fmt.Sprintf("%s.zip", file), &buf)
	if err != nil {
		err = fmt.Errorf("%w:%w", errors.ErrUploadFile, err)
		slog.Error("upload file", "err", err)
		err2 := msgrSvc.SendMessage("framer-status.fifo", fmt.Sprintf("%s.FALHA.%s", file, err))
		if err2 != nil {
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
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
			err2 = fmt.Errorf("%w: %w", errors.ErrSendMessage, err2)
			slog.Error("send message", "err", err2)
		}
		return err
	}

	return nil
}
