package utils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	DBService "parse-log/db"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	IsAverageTime = true
	NoAverageTime = false
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

func ExportCSV(result []bson.M, filePath string, isAverageTime bool) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, value := range result {
		valueStr := (fmt.Sprintf("%s, ", value["_id"]))
		if isAverageTime {
			if _, err := w.WriteString(fmt.Sprintf("%s, %d, %d, %d\n", valueStr, value["proxy"], value["gateway"], value["request"])); err != nil {
				logrus.Fatal("Failed to write on the file")
			}
		} else {
			if _, err := w.WriteString(fmt.Sprintf("%s%d\n", valueStr, value["requests"])); err != nil {
				logrus.Fatal("Failed to write on the file")
			}
		}
	}

	w.Flush()

	return nil
}
