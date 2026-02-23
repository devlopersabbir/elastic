package config

import (
	"fmt"
	"os"
)

func LoadEnv() *Config {
	return &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPass:           os.Getenv("DB_PASS"),
		DBName:           os.Getenv("DB_NAME"),
		ElasticsearchUrl: os.Getenv("ELASTICSEARCH_URL"),
		ServerPort:       os.Getenv("SERVER_PORT"),
	}
}

func (c *Config) PostgresDNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
}

func GetEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
