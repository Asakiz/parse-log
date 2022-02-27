package main

import (
	"context"
	"os"
	DBService "parse-log/db"
	"parse-log/utils"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	if err := utils.CheckFilePath(os.Args); err != nil {
		logrus.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}

	// try to ping mongoDB before go to the next step
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logrus.WithError(err).Fatal("Failed to ping MongoDB")
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logrus.Fatal(err)
		}
	}()

	db := client.Database("main")
	service := DBService.Service{DB: db.Collection("gateway"), Context: context.TODO()}

	logrus.Warn("Populating the database, please wait...")

	if err := utils.PopulateDB(&service, context.TODO(), os.Args[1]); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Done!")

	consumersID, err := service.GetAllIDs(DBService.Consumers)
	if err != nil {
		logrus.Fatal(err)
	}

	// create a index to make the search more fast since the consult to this ID is heavy
	index := mongo.IndexModel{Keys: bson.M{"authenticated_entity.consumer_id.uuid": 1}}
	if _, err := service.DB.Indexes().CreateOne(context.TODO(), index); err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calculating the requests for the consumers, please wait...")

	result, err := service.CalcRequests(consumersID, DBService.Consumers)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calcutations done, generating the csv file")

	if err := utils.ExportCSV(result, "consumer-request.csv", utils.NoAverageTime); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Done!")

	servicesID, err := service.GetAllIDs(DBService.Services)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calculating the requests for the services, please wait...")

	result, err = service.CalcRequests(servicesID, DBService.Services)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calcutations done, generating the csv file")

	if err := utils.ExportCSV(result, "service-request.csv", utils.NoAverageTime); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Done!")

	logrus.Warn("Calculating the average time for proxy, gateway and request for service. Please wait...")

	result, err = service.CalcAverageTime(servicesID)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calcutations done, generating the csv file")

	if err := utils.ExportCSV(result, "average-time-request.csv", utils.IsAverageTime); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Done!")
}
