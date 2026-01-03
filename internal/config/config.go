package config

import "os"

type Config struct {
	Port string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}
}

// Вспомогательная функция для чтения переменных окружения
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
