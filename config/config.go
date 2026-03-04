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
		DBHost:            GetEnv("DB_HOST", "carma.postgres.database.azure.com"),
		DBPort:            GetEnv("DB_PORT", "5432"),
		DBUser:            GetEnv("DB_USER", "carmaadmin"),
		DBPass:            GetEnv("DB_PASS", "Hosthunter1221!."),
		DBName:            GetEnv("DB_NAME", "postgres"),
		DBSSLMode:         sslMode,
		ElasticsearchUrl:  GetEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		ElasticsearchUrls: GetEnv("ELASTICSEARCH_URLS", "http://localhost:9200"), // optional extra nodes
		ServerPort:        GetEnv("SERVER_PORT", "8080"),
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
