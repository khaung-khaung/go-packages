package services

import (
	"context"

	"github.com/banyar/go-packages/pkg/frontlog"
	"github.com/banyar/go-packages/pkg/repositories"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService struct {
	MongoRepo *repositories.MongoRepository
}

func NewMongoService(db *repositories.MongoRepository) *MongoService {
	return &MongoService{
		MongoRepo: db,
	}
}

func (s *MongoService) GetClient() (*mongo.Client, error) {
	client := s.MongoRepo.Client
	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		frontlog.Logger.Error("Error check connection to mongo client :", zap.Any("error=", err))
		return nil, err
	}
	return client, nil
}

func (s *MongoService) GetCollection(col string) (*mongo.Collection, error) {
	collection := s.MongoRepo.Client.Database(s.MongoRepo.DSNMongo.Database).Collection(col)
	return collection, nil
}
