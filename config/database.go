package config

import (
	"log"

	"github.com/devloeprsabbir/go-elasticsearch/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(config *Config) {
	db, err := gorm.Open(postgres.Open(config.PostgresDNS()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	db.AutoMigrate(&models.Todo{})
	DB = db
}
