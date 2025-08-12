package repositories

import (
	"context"
	"fmt"
	"log"

	entities "github.com/banyar/go-packages/pkg/entities"

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
		fmt.Printf("Error connecting to mongo database : error=%v", err)
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Printf("Error check connection to mongo client : error=%v", err)
		log.Fatal(err)
	}

	return &MongoRepository{
		Client:   client,
		DSNMongo: DSNMongo,
	}

}
