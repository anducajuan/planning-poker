package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	Port        string
	LogMode     string
	CORSOrigins string
}

// LoadConfig carrega as configurações das variáveis de ambiente
func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:admin@localhost:5432/database?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		LogMode:     getEnv("LOG_MODE", "dev"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),
	}
}

// getEnv retorna o valor da variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retorna o valor da variável de ambiente como inteiro ou um valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// getEnvAsBool retorna o valor da variável de ambiente como booleano ou um valor padrão
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
