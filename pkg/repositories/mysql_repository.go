package repositories

import (
	"fmt"

	entities "github.com/banyar/go-packages/pkg/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlRepository struct {
	Connection *gorm.DB
}

func ConnectDatabase(DSNMySQL *entities.DSNMySQL) *MysqlRepository {
	DSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		DSNMySQL.Username,
		DSNMySQL.Password,
		DSNMySQL.Host,
		DSNMySQL.Port,
		DSNMySQL.Database,
	)
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
		//panic("Can't connect to DB!")
		db.Error = err
	}

	return &MysqlRepository{
		Connection: db,
	}
}
