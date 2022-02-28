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

const (
	IsAverageTime = true
)

// Function to check if the input file is valid
// If encounters any errors, it will return
func CheckFilePath(filePath []string) error {
	// check if the size of std.input is greater than 2
	if len(filePath) < 2 {
		return ErrMissingArguments
	}

	// match the extension if a regex to check if is valid
	if match, _ := filepath.Match(".txt", filepath.Ext(filePath[1])); !match {
		return ErrFileExtension
	}

	return nil
}

func InitDatabase() (*DBService.Service, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	if err != nil {
		return nil, err
	}

	// try to ping mongoDB before go to the next step
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
	// check if the filePath exists
	file, err := os.Open(filePath)
	if err != nil {
		return ErrToOpenFile
	}
	defer file.Close()

	// create a new buffer reader to the file
	reader := bufio.NewReader(file)
	for {
		var buffer bytes.Buffer
		var l []byte
		var isPrefix bool

		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// check if the line is too long for the buffer, if it's isPrefix is going to be set true
			// then return the beginning of the line, the rest of the line will be returned later
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

		// if encounter a EOF the function stops and return the controller to main
		if err == io.EOF {
			break
		}

		// insert the bytes read into the database
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
func ExportCSV(result []bson.M, filePath string, isAverageTime ...bool) error {
	// check if is possible to create the file with on the filePath
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, value := range result {
		// this line is necessary as it will be repeated in both cases
		valueStr := (fmt.Sprintf("%s, ", value["_id"]))
		if len(isAverageTime) > 0 && isAverageTime[0] {
			// this partically is created because the average time has three values to write
			if _, err := w.WriteString(fmt.Sprintf("%s%d, %d, %d\n", valueStr, value["proxy"], value["gateway"], value["request"])); err != nil {
				return ErrWriteToFile
			}
		} else {
			if _, err := w.WriteString(fmt.Sprintf("%s%d\n", valueStr, value["requests"])); err != nil {
				return ErrWriteToFile
			}
		}
	}

	w.Flush()

	return nil
}
