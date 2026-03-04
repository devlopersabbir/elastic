package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(config *Config) {
	db, err := gorm.Open(postgres.Open(config.PostgresDNS()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// We do NOT call AutoMigrate here — the vehicle_data table already exists in PG.
	// We read from it in read-only mode through the vehicle_marketplace schema.
	log.Println("Connected to PostgreSQL successfully")
	DB = db
}
