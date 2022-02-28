package database

import (
	"context"
	"encoding/json"
	"parse-log/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Context    context.Context
}

type Arguments struct {
	FieldName     string
	IsAverageTime bool
}

// Function to insert the log on the database
// recieve a array of bytes and convert to the gateway struct
func (s *Service) InsertLog(input []byte) error {
	var gateway models.Gateway

	// convert the input into the Gateway model
	if err := json.Unmarshal(input, &gateway); err != nil {
		return err
	}

	// insert on the database a Gateway entry
	if _, err := s.Collection.InsertOne(s.Context, &gateway); err != nil {
		return err
	}

	return nil
}

// Function to get all distinct IDs stored on the database
func (s *Service) GetAllIDs(arg Arguments) ([]interface{}, error) {
	list, err := s.Collection.Distinct(s.Context, string(arg.FieldName), bson.M{})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Function to calculate all the requests based on the input list
func (s *Service) CalculateQuery(list []interface{}, arg Arguments, query bson.M) ([]bson.M, error) {
	var result []bson.M

	for _, id := range list {
		// the logic of the query is to search for all the entries correspond to the ID
		// then sum all the occurrences
		cursor, err := s.Collection.Aggregate(s.Context, []bson.M{
			{"$match": bson.M{string(arg.FieldName): id}}, query})
		if err != nil {
			return nil, err
		}

		// extract the result into a map
		result, err = extractResult(cursor, result)
		if err != nil {
			return nil, err
		}
	}

	if arg.IsAverageTime {
		// this is to calculate the average time of every field
		for _, value := range result {
			value["proxy"] = value["proxy"].(int32) / value["total"].(int32)
			value["kong"] = value["gateway"].(int32) / value["total"].(int32)
			value["request"] = value["request"].(int32) / value["total"].(int32)
		}
	}

	return result, nil
}

// Function to extract the result into a map, so is can be accessed more easily
func extractResult(cursor *mongo.Cursor, result []bson.M) ([]bson.M, error) {
	var showsLoaded []bson.M

	if err := cursor.All(context.TODO(), &showsLoaded); err != nil {
		return nil, err
	}

	result = append(result, showsLoaded...)

	return result, nil
}
