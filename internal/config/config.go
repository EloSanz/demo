package config

import (
	"os"
)

// Config holds all application configuration.
type Config struct {
	Port     string
	DBEngine string // "sqlite" or "postgres"
	DBDSN    string // Data Source Name (path for sqlite, connection string for postgres)
	
	// External APIs
	JSONPlaceholderURL string
	RickAndMortyURL     string

	// AWS (optional tags for documentation/reference)
	S3Bucket string
}

// Load populates Config from environment variables with sensible defaults.
func Load() Config {
	return Config{
		Port:               getEnv("PORT", "8080"),
		DBEngine:           getEnv("DB_ENGINE", "sqlite"),
		DBDSN:              getEnv("DATABASE_URL", "demo.db"),
		JSONPlaceholderURL: getEnv("JSONPLACEHOLDER_BASE_URL", "https://jsonplaceholder.typicode.com"),
		RickAndMortyURL:     getEnv("RICKANDMORTY_BASE_URL", "https://rickandmortyapi.com/api"),
		S3Bucket:           getEnv("AWS_S3_BUCKET", "myawsbucketelito"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
