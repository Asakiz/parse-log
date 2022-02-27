package database

import (
	"context"
	"encoding/json"
	"parse-log/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	DB      *mongo.Collection
	Context context.Context
}

type Arguments string

const (
	Consumers Arguments = "authenticated_entity.consumer_id.uuid"
	Services  Arguments = "service.id"
)

// Function to insert the log on the database
// recieve a array of bytes and convert to the gateway struct
func (s *Service) InsertLog(input []byte) error {
	var gateway models.Gateway

	// convert the input into the Gateway model
	if err := json.Unmarshal(input, &gateway); err != nil {
		return err
	}

	// insert on the database a Gateway entry
	if _, err := s.DB.InsertOne(s.Context, &gateway); err != nil {
		return err
	}

	return nil
}

// Function to get all distinct IDs stored on the database
func (s *Service) GetAllIDs(arg Arguments) ([]interface{}, error) {
	list, err := s.DB.Distinct(s.Context, string(arg), bson.M{})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Function to calculate all the requests based on the input list
func (s *Service) CalcRequests(List []interface{}, arg Arguments) ([]bson.M, error) {
	var result []bson.M

	for _, id := range List {
		// the logic of the query is to search for all the entries correspond to the ID
		// then sum all the occurrences
		cursor, err := s.DB.Aggregate(s.Context, []bson.M{
			{"$match": bson.M{string(arg): id}},
			{"$group": bson.M{"_id": "$" + string(arg), "requests": bson.M{"$sum": 1}}}})
		if err != nil {
			return nil, err
		}

		// extract the result into a map
		result, err = extractResult(cursor, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Function to calculate the average time of all the latencies
func (s *Service) CalcAverageTime(List []interface{}) ([]bson.M, error) {
	var result []bson.M

	for _, id := range List {
		// the logic of the query is to search for all the entries correspond to the ID
		// then sum all the occurrences and the fields: proxy, gateway and request
		cursor, err := s.DB.Aggregate(s.Context, []bson.M{
			{"$match": bson.M{"service.id": id}},
			{"$group": bson.M{"_id": "$service.id", "proxy": bson.M{"$sum": "$latencies.proxy"},
				"gateway": bson.M{"$sum": "$latencies.gateway"},
				"request": bson.M{"$sum": "$latencies.request"},
				"total":   bson.M{"$sum": 1},
			}}})
		if err != nil {
			return nil, err
		}

		// extract the result into a map
		result, err = extractResult(cursor, result)
		if err != nil {
			return nil, err
		}
	}

	// this is to calculate the average time of every field
	for _, value := range result {
		value["proxy"] = value["proxy"].(int32) / value["total"].(int32)
		value["kong"] = value["gateway"].(int32) / value["total"].(int32)
		value["request"] = value["request"].(int32) / value["total"].(int32)
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
