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

func (s *Service) InsertLog(input []byte) error {
	var gateway models.Gateway

	if err := json.Unmarshal(input, &gateway); err != nil {
		return err
	}

	if _, err := s.DB.InsertOne(s.Context, &gateway); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAllIDs(arg Arguments) ([]interface{}, error) {
	consumerList, err := s.DB.Distinct(s.Context, string(arg), bson.M{})
	if err != nil {
		return nil, err
	}

	return consumerList, nil
}

func (s *Service) CalcRequests(List []interface{}, arg Arguments) []bson.M {
	var result []bson.M

	for _, id := range List {
		showLoadedCursor, err := s.DB.Aggregate(s.Context, []bson.M{
			{"$match": bson.M{string(arg): id}},
			{"$group": bson.M{"_id": "$" + string(arg), "requests": bson.M{"$sum": 1}}}})
		if err != nil {
			panic(err)
		}

		var showsLoaded []bson.M
		if err = showLoadedCursor.All(s.Context, &showsLoaded); err != nil {
			panic(err)
		}

		result = append(result, showsLoaded...)
	}

	return result
}

func (s *Service) CalcAverageTime(List []interface{}) []bson.M {
	var result []bson.M

	for _, id := range List {
		showLoadedCursor, err := s.DB.Aggregate(s.Context, []bson.M{
			{"$match": bson.M{"service.id": id}},
			{"$group": bson.M{"_id": "$service.id", "proxy": bson.M{"$sum": "$latencies.proxy"},
				"gateway": bson.M{"$sum": "$latencies.gateway"},
				"request": bson.M{"$sum": "$latencies.request"},
			}}})
		if err != nil {
			panic(err)
		}

		var showsLoaded []bson.M
		if err = showLoadedCursor.All(s.Context, &showsLoaded); err != nil {
			panic(err)
		}

		result = append(result, showsLoaded...)
	}

	return result
}
