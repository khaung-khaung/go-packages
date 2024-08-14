package services

import (
	"fmt"

	"github.com/banyar/go-packages/pkg/repositories"

	"gorm.io/gorm"
)

type MysqlService struct {
	DbRepo *repositories.MysqlRepository
}

func NewMysqlService(db *repositories.MysqlRepository) *MysqlService {
	return &MysqlService{
		DbRepo: db,
	}
}

func (s *MysqlService) GetConnection() (*gorm.DB, error) {
	db := s.DbRepo.Connection
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
		return nil, err
	}
	sqlDB.SetMaxOpenConns(5)
	return db, nil
}
