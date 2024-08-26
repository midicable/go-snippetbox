package config

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}
}

func New() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "root"),
		Password: getEnv("POSTGRES_PASSWORD", "password"),
		DBName:   getEnv("POSTGRES_DB", "postgres"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exist := os.LookupEnv(key)

	if exist {
		return value
	}

	return defaultValue
}
