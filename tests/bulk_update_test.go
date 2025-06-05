package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestBulkUpdateFailures(t *testing.T) {

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
	collection, err := client.MongoService.GetCollection("node")
	fmt.Println("collection  ", collection)

	if err != nil {
		fmt.Println("error  ", err)
	}

	var existingTags []string
	existingTags = append(existingTags, "str")

	objectID, err := primitive.ObjectIDFromHex("5f8db962890fab91e6467a5e")
	if err != nil {
		fmt.Println("error  ", err)
	}

	update := bson.M{"$set": bson.M{"node_attr.tags": existingTags}}
	tests := []struct {
		name          string
		writeModels   []mongo.WriteModel
		expectedError bool
	}{
		{
			name: "Network issue simulation",
			writeModels: []mongo.WriteModel{
				mongo.NewUpdateOneModel().
					SetFilter(bson.M{"_id": objectID}).
					SetUpdate(update),
			},
			expectedError: true,
		},
		// {
		// 	name: "Incorrect query syntax",
		// 	writeModels: []mongo.WriteModel{
		// 		mongo.NewUpdateOneModel().
		// 			SetFilter(bson.M{"_id": "invalid"}). // Intentional mistake
		// 			SetUpdate(update),
		// 	},
		// 	expectedError: true,
		// },
		// {
		// 	name: "Schema violation",
		// 	writeModels: []mongo.WriteModel{
		// 		mongo.NewUpdateOneModel().
		// 			SetFilter(bson.M{"_id": objectID1}).
		// 			SetUpdate(update),
		// 	},
		// 	expectedError: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the bulk write
			bulkOptions := options.BulkWrite().SetOrdered(false)
			result, err := collection.BulkWrite(context.TODO(), tt.writeModels, bulkOptions)
			if (err != nil) != tt.expectedError {
				t.Errorf("BulkWrite() error = %v, expectedError %v", err, tt.expectedError)
			}
			fmt.Println("tests  ", result.MatchedCount, result.ModifiedCount)
		})
	}
}
