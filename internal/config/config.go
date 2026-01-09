package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppPort                string
	PostgresDSN            string
	APIKey                 string
	StatsTimeWindowMinutes int
	WebhookURL             string
}

func Load() *Config {
	appPort := getEnv("APP_PORT", "8080")
	postgresDSN := getEnv("POSTGRES_DSN", "")
	apiKey := getEnv("API_KEY", "secret123")
	statsMinutesStr := getEnv("STATS_TIME_WINDOW_MINUTES", "5")
	webhookURL := getEnv("WEBHOOK_URL", "")

	statsMinutes, err := strconv.Atoi(statsMinutesStr)
	if err != nil {
		log.Fatal("invalid STATS_TIME_WINDOW_MINUTES")
	}

	if postgresDSN == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	if apiKey == "" {
		log.Fatal("API_KEY is required")
	}

	return &Config{
		AppPort:                appPort,
		PostgresDSN:            postgresDSN,
		APIKey:                 apiKey,
		StatsTimeWindowMinutes: statsMinutes,
		WebhookURL:             webhookURL,
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
