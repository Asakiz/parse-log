package main

import (
	"context"
	"os"
	"parse-log/utils"
	"time"

	"github.com/sirupsen/logrus"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	if err := utils.CheckFilePath(os.Args); err != nil {
		logrus.Fatal(err)
	}

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
			logrus.Fatal(err)
		}
	}()

	logrus.Warn("Populating the Database, please wait...")

	if err := utils.PopulateDB(client, ctx, os.Args[1]); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Done!")
}
