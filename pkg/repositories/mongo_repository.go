package repositories

import (
	"context"
	"fmt"

	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/frontlog"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client   *mongo.Client
	DSNMongo *entities.DSNMongo
}

func ConnectMongo(DSNMongo *entities.DSNMongo) *MongoRepository {

	// MongoDB connection URI with authentication
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", DSNMongo.Username, DSNMongo.Password, DSNMongo.Host, DSNMongo.Port, DSNMongo.Database)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		frontlog.Logger.Error(
			"Error check connection to mongo database:",
			zap.Any("err", err.Error()),
		)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		frontlog.Logger.Error(
			"Error check connection to mongo client:",
			zap.Any("err", err.Error()),
		)
	}

	return &MongoRepository{
		Client:   client,
		DSNMongo: DSNMongo,
	}

}
