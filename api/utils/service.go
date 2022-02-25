package utils

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	DBService "parse-log/db"
)

func CheckFilePath(filePath []string) error {
	if len(filePath) < 2 {
		return ErrMissingArguments
	}

	if match, _ := filepath.Match(".txt", filepath.Ext(filePath[1])); !match {
		return ErrFileExtension
	}

	return nil
}

func PopulateDB(service *DBService.Service, ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		var buffer bytes.Buffer
		var l []byte
		var isPrefix bool

		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			if !isPrefix {
				break
			}

			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
		}

		if err == io.EOF {
			break
		}

		if err := service.InsertLog(buffer.Bytes()); err != nil {
			return err
		}
	}

	if err != io.EOF {
		return err
	}

	return nil
}
