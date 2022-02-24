package utils

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	DBService "parse-log/db"

	mongo "go.mongodb.org/mongo-driver/mongo"
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

func PopulateDB(client *mongo.Client, ctx context.Context, filePath string) error {
	db := client.Database("main")
	service := DBService.Service{DB: db, Context: ctx}

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

		service.InsertLog(buffer.Bytes())

		if err == io.EOF {
			break
		}
	}
	if err != io.EOF {
		return err
	}

	return nil
}
