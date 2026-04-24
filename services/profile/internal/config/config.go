package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                    string
	LogLevel                string
	DatabaseURL             string
	JWTSecret               string
	RedpandaBrokers         []string
	RedpandaOrderPaidTopic  string
	RedpandaConsumerGroupID string
}

func Load() *Config {
	return &Config{
		Port:                    getEnv("PORT", "8003"),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		DatabaseURL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/litsee_profile?sslmode=disable"),
		JWTSecret:               getEnv("JWT_SECRET", "hihihaha"),
		RedpandaBrokers:         splitCSV(getEnv("REDPANDA_BROKERS", "localhost:9092")),
		RedpandaOrderPaidTopic:  getEnv("REDPANDA_ORDER_PAID_TOPIC", "orderPaid"),
		RedpandaConsumerGroupID: getEnv("REDPANDA_CONSUMER_GROUP", "profile-service"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
