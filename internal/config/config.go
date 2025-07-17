// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func parseEnv() error {

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Println("godotenv.Load().Error:", err)
		return nil
	}

	missing := make([]string, 0)
	envVars := []string{}

	for _, v := range envVars {
		envVal := os.Getenv(v)
		if envVal == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	return fmt.Errorf("missing env vars: %v", missing)
}

func LoadAppConfig() *Config {
	env := parseEnv()
	if env != nil {
		log.Fatal("Failed to load app config:", env)
	}

	jwtExpiry, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRY"), 10, 64)

	if jwtExpiry == 0 {
		jwtExpiry = 1 // Default to 1 hour
	}

	jwtRefreshExpiry, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRY"), 10, 64)
	if jwtRefreshExpiry == 0 {
		jwtRefreshExpiry = 7 // Default to 1 week
	}

	config := Config{
		KAFKA_BROKER_URL:      os.Getenv("KAFKA_BROKER_URL"),
		PostgresDB:            os.Getenv("POSTGRES_DB"),
		PostgresUser:          os.Getenv("POSTGRES_USER"),
		PostgresPass:          os.Getenv("POSTGRES_PASSWORD"),
		PostgresDSN:           os.Getenv("POSTGRES_DSN"),
		RedisAddr:             os.Getenv("REDIS_ADDR"),
		RedisPassword:         os.Getenv("REDIS_PASSWORD"),
		PORT:                  os.Getenv("PORT"),
		ORIGINS:               os.Getenv("ORIGINS"),
		PostgresHost:          os.Getenv("POSTGRES_HOST"),
		PostgresPort:          os.Getenv("POSTGRES_PORT"),
		JwtSecret:             os.Getenv("JWT_SECRET"),
		RefreshTokenSecret:    os.Getenv("REFRESH_TOKEN_SECRET"),
		EnableDBQueryLogging:  os.Getenv("ENABLE_DB_QUERY_LOGGING") == "true",
		JwtTokenRefreshExpiry: jwtRefreshExpiry,
		JwtTokenExpiry:        jwtExpiry,
	}

	return &config
}
