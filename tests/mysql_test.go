package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/entities"
)

func TestMysqlDatabase(t *testing.T) {
	DSNMySQL := entities.DSNMySQL{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_POST"),
		Username: os.Getenv("MYSQL_USERNAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}

	databaseAdapter := adapters.NewDatabaseAdapter(&DSNMySQL)

	connection, err := databaseAdapter.DatabaseService.GetConnection()
	if err != nil {
		log.Fatal("ERROR : ", err)
	}

	fmt.Println("myswl database Adapter ===> ", connection)
}
