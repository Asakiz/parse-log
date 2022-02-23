package database

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"parse-log/models"
)

type Service struct {
	DB      *mongo.Database
	Context context.Context
}

func (s *Service) InsertLog(input []byte) error {
	var gateway models.Gateway

	if err := json.Unmarshal(input, &gateway); err != nil {
		return err
	}

	if _, err := s.DB.Collection("gateway").InsertOne(s.Context, gateway); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteLog(clientIP string) error {
	if _, err := s.DB.Collection("gateway").DeleteOne(s.Context, bson.M{"clientip": clientIP}); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetLog(clientIP string) (*models.Gateway, error) {
	var gateway models.Gateway

	if err := s.DB.Collection("gateway").FindOne(s.Context, bson.M{"clientip": clientIP}).Decode(&gateway); err != nil {
		return nil, err
	}

	return &gateway, nil
}

func (s *Service) GetRequestsByConsumer(clientIP string) uint16 {
	gateway, err := s.GetLog(clientIP)

	if err != nil {
		return 0
	}

	return gateway.Request.Size
}

func (s *Service) GetRequestsByService(clientIP string) uint16 {
	gateway, err := s.GetLog(clientIP)

	if err != nil {
		return 0
	}

	return gateway.Service.Retries
}
