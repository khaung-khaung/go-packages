package services

import (
	"context"
	"fmt"
	"log"

	"github.com/banyar/go-packages/pkg/repositories"

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
		fmt.Printf("Error check connection to mongo client : error=%v", err)
		log.Fatal(err)
		return nil, err
	}
	return client, nil
}

func (s *MongoService) GetCollection(col string) (*mongo.Collection, error) {
	collection := s.MongoRepo.Client.Database(s.MongoRepo.DSNMongo.Database).Collection(col)
	return collection, nil
}
