package config

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	DBSSLMode         string
	ElasticsearchUrl  string
	ElasticsearchUrls string // comma-separated list for multiple nodes
	ServerPort        string
}
