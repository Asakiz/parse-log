package main

import (
	"context"
	"os"
	DBService "parse-log/db"
	"parse-log/utils"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if err := utils.CheckFilePath(os.Args); err != nil {
		logrus.Fatal(err)
	}

	db, err := utils.InitDatabase()
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		if err := db.Client.Disconnect(context.TODO()); err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Warn("Populating the database, please wait...")

	if err := utils.PopulateDB(db, context.TODO(), os.Args[1]); err != nil {
		logrus.Fatal(err)
	}

	consumersID, err := db.GetAllIDs(DBService.Arguments{FieldName: "authenticated_entity.consumer_id.uuid"})
	if err != nil {
		logrus.Fatal(err)
	}

	index := mongo.IndexModel{Keys: bson.M{"authenticated_entity.consumer_id.uuid": 1}}
	if _, err := db.Collection.Indexes().CreateOne(context.TODO(), index); err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calculating the requests for the consumers")

	result, err := db.CalculateQuery(consumersID,
		DBService.Arguments{FieldName: "authenticated_entity.consumer_id.uuid"},
		bson.M{"$group": bson.M{"_id": "$authenticated_entity.consumer_id.uuid", "requests": bson.M{"$sum": 1}}})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calcutations done, generating the csv file")

	if err := utils.ExportCSV(utils.Output{FilePath: "consumer-request.csv",
		Result: utils.ResultMap{
			Value: result, Field: []string{"requests"}}}); err != nil {
		logrus.Fatal(err)
	}

	servicesID, err := db.GetAllIDs(DBService.Arguments{FieldName: "service.id"})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calculating the requests for the services")

	result, err = db.CalculateQuery(servicesID,
		DBService.Arguments{FieldName: "service.id"},
		bson.M{"$group": bson.M{"_id": "$service.id", "requests": bson.M{"$sum": 1}}})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calcutations done, generating the csv file")

	if err := utils.ExportCSV(utils.Output{FilePath: "service-request.csv",
		Result: utils.ResultMap{Value: result, Field: []string{"requests"}}}); err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calculating the average time for proxy, gateway and request for service")

	result, err = db.CalculateQuery(servicesID,
		DBService.Arguments{FieldName: "service.id", IsAverageTime: true},
		bson.M{"$group": bson.M{"_id": "$service.id",
			"proxy":   bson.M{"$sum": "$latencies.proxy"},
			"gateway": bson.M{"$sum": "$latencies.gateway"},
			"request": bson.M{"$sum": "$latencies.request"},
			"total":   bson.M{"$sum": 1}}})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Warn("Calcutations done, generating the csv file")

	if err := utils.ExportCSV(utils.Output{
		FilePath: "average-time-request.csv",
		Result: utils.ResultMap{
			Value: result,
			Field: []string{
				"proxy",
				"gateway",
				"request",
			},
		}}); err != nil {
		logrus.Fatal(err)
	}
}
