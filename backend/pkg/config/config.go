package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {

	// Server
	Port string
	Env  string

	// Database
	MongoURI      string
	MongoDatabase string

	// TwelveData API
	TwelveDataAPIKey string

	// Finnhub API
	FinnhubAPIKey string

	// JWT
	JWTSecret         string
	JWTExpireDuration time.Duration

	// Redis (Optional)
	RedisURI string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		TwelveDataAPIKey: getEnv("TWELVEDATA_API_KEY", ""),
		FinnhubAPIKey:    getEnv("FINNHUB_API_KEY", ""),

		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "bengi-investment-system"),

		JWTSecret:         getEnv("JWT_SECRET", "secret"),
		JWTExpireDuration: parseDuration(getEnv("JWT_EXPIRE_DURATION", "24h")),

		RedisURI: getEnv("REDIS_URI", "redis://localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
