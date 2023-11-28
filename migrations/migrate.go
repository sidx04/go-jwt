package migrations

import (
	"github.com/sidx04/go-jwt/initialisers"
	"github.com/sidx04/go-jwt/models"
)

func MigrateDatabase() {
	initialisers.DB.AutoMigrate(&models.User{})
}
