package postgres

import (
	"fmt"
	"log"

	"github.com/amir79esmaeili/go-tel-money/internal/cfg"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const UrlSample = "postgres://%s:%s@%s:%s/%s"

// ConnectToDB takes configuration struct and creates a connection to the given Database
func ConnectToDB(cfg *cfg.Config) (*gorm.DB, error) {
	dbUrl := fmt.Sprintf(UrlSample, cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.ExpenseType{})
	db.AutoMigrate(&models.Expense{})
	log.Println("[INFO] Database connected...")
	return db, nil
}
