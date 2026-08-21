package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost  string
	DBPort  string
	DBUser  string
	DBPwrd  string
	DBName  string
	AppPort string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	return &Config{
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBPort:  getEnv("DB_PORT", "5432"),
		DBUser:  getEnv("DB_USER", "postgres"),
		DBPwrd:  getEnv("DB_PWRD", "postgres"),
		DBName:  getEnv("DB_NAME", "finance"),
		AppPort: getEnv("APP_PORT", "3000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
