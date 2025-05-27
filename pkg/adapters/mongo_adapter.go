package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"
)

type MongoAdapter struct {
	MongoService interfaces.IMongoService
}

func NewMongoAdapter(DSNMongo *entities.DSNMongo) *MongoAdapter {
	databaseRepo := repositories.ConnectMongo(DSNMongo)
	return &MongoAdapter{
		MongoService: services.NewMongoService(databaseRepo),
	}
}
