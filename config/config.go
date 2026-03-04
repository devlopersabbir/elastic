package config

import (
	"fmt"
	"os"
	"strings"
)

func LoadEnv() *Config {
	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "require" // Azure PG always requires SSL
	}
	return &Config{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPass:            os.Getenv("DB_PASS"),
		DBName:            os.Getenv("DB_NAME"),
		DBSSLMode:         sslMode,
		ElasticsearchUrl:  os.Getenv("ELASTICSEARCH_URL"),
		ElasticsearchUrls: os.Getenv("ELASTICSEARCH_URLS"), // optional extra nodes
		ServerPort:        os.Getenv("SERVER_PORT"),
	}
}

func (c *Config) PostgresDNS() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}

// ElasticAddresses returns the list of Elasticsearch addresses to connect to.
// Falls back to ElasticsearchUrl if no multi-node list is provided.
func (c *Config) ElasticAddresses() []string {
	if c.ElasticsearchUrls != "" {
		var addrs []string
		for _, u := range strings.Split(c.ElasticsearchUrls, ",") {
			u = strings.TrimSpace(u)
			if u != "" {
				addrs = append(addrs, u)
			}
		}
		if len(addrs) > 0 {
			return addrs
		}
	}
	return []string{c.ElasticsearchUrl}
}

func GetEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
