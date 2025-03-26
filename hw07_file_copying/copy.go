package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			log.Printf("Ошибка закрытия файла: %v", err)
		}
	}()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// проверяем, что это файл и он не пустой
	if fileInfo.Size() == 0 || fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// фиксируем limit в случае полного чтения файла
	if limit == 0 || limit > fileInfo.Size()-offset {
		limit = fileInfo.Size() - offset
	}

	// информация для отладки
	fmt.Printf("Входной файл: размер %d, offset %d, limit %d\n", fileInfo.Size(), offset, limit)

	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		err := dstFile.Close()
		if err != nil {
			log.Printf("Ошибка закрытия файла: %v", err)
		}
	}()

	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	reader := bar.NewProxyReader(srcFile)

	_, err = io.CopyN(dstFile, reader, limit)
	if err != nil {
		return err
	}

	return nil
}
