package config

import (
	"os"
)

// Config holds all application configuration.
type Config struct {
	Port     string
	// Engines
	DBEngine           string
	NotificationEngine string // "sqs" or "memory"
	DBDSN    string // Data Source Name (path for sqlite, connection string for postgres)
	
	// External APIs
	JSONPlaceholderURL string
	RickAndMortyURL     string
	ElasticsearchURL    string

	// AWS (optional tags for documentation/reference)
	S3Bucket           string
	SQSQueueURL        string
}

// Load populates Config from environment variables with sensible defaults.
func Load() Config {
	return Config{
		Port:               getEnv("PORT", "8080"),
		DBEngine:           getEnv("DB_ENGINE", "sqlite"),
		NotificationEngine: getEnv("NOTIFICATION_ENGINE", "memory"),
		DBDSN:              getEnv("DATABASE_URL", "demo.db"),
		JSONPlaceholderURL: getEnv("JSONPLACEHOLDER_BASE_URL", "https://jsonplaceholder.typicode.com"),
		RickAndMortyURL:     getEnv("RICKANDMORTY_BASE_URL", "https://rickandmortyapi.com/api"),
		ElasticsearchURL:    getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		S3Bucket:           getEnv("AWS_S3_BUCKET", "myawsbucketelito"),
		SQSQueueURL:        getEnv("AWS_SQS_QUEUE_URL", "https://sqs.us-east-1.amazonaws.com/123456789/my-queue"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
