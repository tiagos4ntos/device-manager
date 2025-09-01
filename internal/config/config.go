package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	ServerPort  string
	HttpTimeout int

	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

var (
	cfg  *Config
	once sync.Once

	DefaultAppName              = "device-manager"
	DefaultPostgresHost         = "localhost"
	DefaultPostgresPort         = "5432"
	DefaultPostgresUserName     = "postgres"
	DefaultPostgresPassword     = "password"
	DefaultPostgresDatabaseName = "device_manager"
)

func LoadConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, falling back to system env")
		}

		cfg = &Config{
			AppName:     getEnvOrDefaultValue("APP_NAME", DefaultAppName),
			ServerPort:  getEnvOrDefaultValue("SERVER_PORT", "8080"),
			HttpTimeout: getIntFromValue(getEnvOrDefaultValue("HTTP_TIMEOUT_IN_SECONDS", "10")),

			DatabaseHost: getEnvOrDefaultValue("DATABASE_HOST", DefaultPostgresHost),
			DatabasePort: getEnvOrDefaultValue("DATABASE_PORT", DefaultPostgresPort),
			DatabaseUser: getEnvOrDefaultValue("DATABASE_USER", DefaultPostgresUserName),
			DatabasePass: os.Getenv("DATABASE_PASS"),
			DatabaseName: getEnvOrDefaultValue("DATABASE_NAME", DefaultPostgresDatabaseName),
		}
	})
	return cfg
}

func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return errors.New("server port is required")
	}
	if c.DatabaseHost == "" {
		return errors.New("database host is required")
	}
	if c.DatabasePort == "" {
		return errors.New("database port is required")
	}
	if c.DatabaseUser == "" {
		return errors.New("database user is required")
	}
	if c.DatabasePass == "" {
		return errors.New("database password is required")
	}
	if c.DatabaseName == "" {
		return errors.New("database name is required")
	}
	return nil
}

func getEnvOrDefaultValue(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntFromValue(value string) int {
	intValue, err := strconv.Atoi(value)

	if err != nil {
		return 0
	}
	return intValue
}
