package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"
	"github.com/banyar/go-packages/pkg/config"
	"github.com/banyar/go-packages/pkg/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestConfig(t *testing.T) {
	fmt.Println("Applicatin Name", config.AppName)
}

func TestMongo(t *testing.T) {
	mongoPort, err := strconv.Atoi(os.Getenv("MONGO_POST"))
	if err != nil {
		log.Fatalf("Error converting MONGO_PORT to integer: %v", err)
	}
	DSNMongo := entities.DSNMongo{
		Host:     os.Getenv("MONGO_HOST"),
		Port:     mongoPort,
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Database: os.Getenv("MONGO_DATABASE"),
	}

	// print("End Adding the User Roles.");
	client := adapters.NewMongoAdapter(&DSNMongo)
	fmt.Println("client Adapter ===> ", client)
	collection, err := client.MongoService.GetCollection("test")
	if err != nil {
		log.Fatal("ERROR : ", err)
	}
	common.DisplayJsonFormat("MongoCollection ===> ", collection)

	filter := bson.D{{Key: "name", Value: "John Doe"}}

	var result bson.M

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("No document was found with the specified filter")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Found a single document:", result)
	}

}
