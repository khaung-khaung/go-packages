package interfaces

import "go.mongodb.org/mongo-driver/mongo"

type IMongoService interface {
	GetClient() (*mongo.Client, error)
	GetCollection(col string) (*mongo.Collection, error)
}
