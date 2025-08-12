package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	entities "github.com/banyar/go-packages/pkg/entities"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client   *mongo.Client
	DSNMongo *entities.DSNMongo
}

func ConnectMongo(DSNMongo *entities.DSNMongo) *MongoRepository {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construct MongoDB connection URI with authSource
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
		DSNMongo.Username,
		DSNMongo.Password,
		DSNMongo.Host,
		DSNMongo.Port,
		DSNMongo.Database,
	)

	// Configure client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// clientOptions := options.Client().
	// 	ApplyURI(mongoURI).
	// 	SetAuth(options.Credential{
	// 		AuthMechanism: "SCRAM-SHA-256",
	// 		Username:      DSNMongo.Username,
	// 		Password:      DSNMongo.Password,
	// 	})

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify connection with ping
	err = client.Ping(ctx, nil)
	if err != nil {
		// Clean up connection if ping fails
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("Failed to disconnect after ping failure: %v", disconnectErr)
		}
		log.Fatalf("Failed to verify MongoDB connection: %v", err)
	}

	// log.Printf("Successfully connected to MongoDB database: %s", DSNMongo.Database)

	return &MongoRepository{
		Client:   client,
		DSNMongo: DSNMongo,
	}
}

// func ConnectMongo(DSNMongo *entities.DSNMongo) *MongoRepository {

// 	// MongoDB connection URI with authentication
// 	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", DSNMongo.Username, DSNMongo.Password, DSNMongo.Host, DSNMongo.Port, DSNMongo.Database)

// 	// Set client options
// 	clientOptions := options.Client().ApplyURI(mongoURI).
// 		SetAuth(options.Credential{
// 			AuthMechanism: "SCRAM-SHA-256", // Use "SCRAM-SHA-1" if needed
// 			Username:      DSNMongo.Username,
// 			Password:      DSNMongo.Password,
// 		})

// 	fmt.Println("clientOptions", clientOptions)
// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		fmt.Printf("Error connecting to mongo database : error=%v", err)
// 		log.Fatal(err)
// 	}

// 	// Check the connection
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		fmt.Printf("Error check connection to mongo client : error=%v", err)
// 		log.Fatal(err)
// 	}

// 	return &MongoRepository{
// 		Client:   client,
// 		DSNMongo: DSNMongo,
// 	}

// }
