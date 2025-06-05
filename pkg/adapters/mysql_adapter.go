package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"
)

type DatabaseAdapter struct {
	DatabaseService interfaces.IDatabaseService
}

func NewDatabaseAdapter(DSNMySQL *entities.DSNMySQL) *DatabaseAdapter {
	databaseRepo := repositories.ConnectDatabase(DSNMySQL)
	return &DatabaseAdapter{
		DatabaseService: services.NewMysqlService(databaseRepo),
	}

}
