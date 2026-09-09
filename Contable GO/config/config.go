package config

import (
	"os"
)

type Config struct {
	Environment    string
	DatabaseDriver string
	DatabaseURL    string
	NodeID         string
}

func LoadFromEnv() *Config {
	driver := getEnv("DB_DRIVER", "sqlite")
	dbURL := getEnv("DATABASE_URL", ":memory:")
	nodeID := getEnv("NODE_ID", "node-local-01")
	env := getEnv("APP_ENV", "development")

	return &Config{
		Environment:    env,
		DatabaseDriver: driver,
		DatabaseURL:    dbURL,
		NodeID:         nodeID,
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}
