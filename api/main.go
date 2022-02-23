package main

import (
	"context"
	"fmt"
	"log"
	"os"
	DBService "parse-log/db"
	"time"

	"github.com/sirupsen/logrus"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	input, err := os.Open("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.WithError(err).Fatal("Failed to ping MongoDB")
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("main")
	service := DBService.Service{DB: db, Context: ctx}
	service.InsertLog(input)
}
