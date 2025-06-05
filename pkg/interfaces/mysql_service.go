package interfaces

import (
	"gorm.io/gorm"
)

type IDatabaseService interface {
	GetConnection() (*gorm.DB, error)
}
