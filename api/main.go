package main

import (
	"context"
	"os"
	DBService "parse-log/db"
	"parse-log/utils"

	"github.com/sirupsen/logrus"

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

	logrus.Warn("Calculating the requests for the consumers, please wait...")

	result := service.CalcRequests(consumersID, DBService.Consumers)

	logrus.Warn("Calcutations done, generating the csv file")

	utils.ExportCSV(result, "output/consumer-request.csv", utils.NoAverageTime)

	logrus.Info("Done!")

	servicesID, err := service.GetAllIDs(DBService.Services)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calculating the requests for the services, please wait...")

	result = service.CalcRequests(servicesID, DBService.Services)

	logrus.Warn("Calcutations done, generating the csv file")

	utils.ExportCSV(result, "output/service-request.csv", utils.NoAverageTime)

	logrus.Info("Done!")

	logrus.Warn("Calculating the average time for proxy, gateway and request for service. Please wait...")

	result = service.CalcAverageTime(servicesID)

	logrus.Warn("Calcutations done, generating the csv file")

	utils.ExportCSV(result, "output/average-time-request.csv", utils.IsAverageTime)

	logrus.Info("Done!")
}
