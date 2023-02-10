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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Output struct {
	FilePath string
	Result   ResultMap
}

type ResultMap struct {
	Value []bson.M
	Field []string
}

// Function to check if the input file is valid
// If encounters any errors, it will return
func CheckFilePath(filePath []string) error {
	if len(filePath) < 2 {
		return ErrMissingArguments
	}

	if match, _ := filepath.Match(".txt", filepath.Ext(filePath[1])); !match {
		return ErrFileExtension
	}

	return nil
}

func InitDatabase() (*DBService.Service, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://parse-log-mongodb:27017/"))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	dbService := DBService.Service{Client: client, Collection: client.Database("main").Collection("gateway"), Context: context.TODO()}

	return &dbService, nil
}

// Function to populate the database base on the input file
// read line by line of the file and insert on database
// util reach the EOF.
func PopulateDB(service *DBService.Service, ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return ErrToOpenFile
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

// function to export the results to a CSV file
// save the results of the calculation based on the filePath argument
func ExportCSV(output Output) error {
	f, err := os.Create(output.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	var fieldStr string

	for _, value := range output.Result.Value {
		valueStr := (fmt.Sprintf("%s, ", value["_id"]))
		if _, err := w.WriteString(valueStr); err != nil {
			return ErrWriteToFile
		}

		for i, field := range output.Result.Field {
			if i-1 == len(output.Result.Field) {
				fieldStr = fmt.Sprintf("%d", value[field])
			} else {
				fieldStr = fmt.Sprintf("%d, ", value[field])
			}

			if _, err := w.WriteString(fieldStr); err != nil {
				return ErrWriteToFile
			}
		}

		w.WriteString("\n")
	}

	w.Flush()

	return nil
}
