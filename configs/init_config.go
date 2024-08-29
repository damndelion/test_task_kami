package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitConfigs() (AppConfigs, error) {
	if os.Getenv("LOCAL") == "true" {
		if err := godotenv.Load(".env.local"); err != nil {
			return AppConfigs{}, fmt.Errorf("error loading .env.local file: %w", err)
		}
	} else {
		if err := godotenv.Load(".env"); err != nil {
			return AppConfigs{}, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	postgres, err := loadPostgresConfig()
	if err != nil {
		return AppConfigs{}, fmt.Errorf("error loading database config: %w", err)
	}

	http, err := loadHttpServerConfig()
	if err != nil {
		return AppConfigs{}, fmt.Errorf("error loading http server config: %w", err)
	}
	logs, err := loadLoggerConfig()
	if err != nil {
		return AppConfigs{}, fmt.Errorf("error loading logger config: %w", err)
	}

	return AppConfigs{
		Postgres: postgres,
		Http:     http,
		Logger:   logs,
	}, nil
}

func loadPostgresConfig() (Postgres, error) {
	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USERNAME", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return Postgres{}, fmt.Errorf("missing required environment variable: %s", v)
		}
	}

	return Postgres{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}, nil
}

func loadHttpServerConfig() (Http, error) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return Http{}, fmt.Errorf("missing required environment variable: %s", "SERVER_PORT")

	}
	return Http{
		Port: port,
	}, nil
}

func loadLoggerConfig() (Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return Logger{}, fmt.Errorf("missing required environment variable: %s", "LOG_LEVEL")

	}
	return Logger{
		LogLevel: logLevel,
	}, nil
}
