package configs

import (
	"os"
	"strconv"
	"github.com/danmaina/logger"
	"github.com/joho/godotenv"
)

// ReadConfigs reads configuration from .env file
func ReadConfigs() (*Config, error) {
	logger.DEBUG("Reading .env config file")

	err := godotenv.Load()
	if err != nil {
		logger.ERR("Error loading .env file, falling back to environment variables", err)
		// We don't return here so it can fallback to existing OS env vars if present
	}

	logLevelStr := os.Getenv("GENERAL.LOG_LEVEL")
	logLevel, _ := strconv.Atoi(logLevelStr)
	if logLevel == 0 {
		logLevel = 3 // default
	}

	config := &Config{
		ApplicationConfigs: GeneralConfigs{
			Port:     os.Getenv("GENERAL.PORT"),
			LogLevel: logLevel,
		},
		Email: EmailConfigs{
			Host:     os.Getenv("EMAIL.HOST"),
			Port:     os.Getenv("EMAIL.PORT"),
			Username: os.Getenv("EMAIL.USERNAME"),
			Password: os.Getenv("EMAIL.PASSWORD"),
		},
		RabbitMQ: RabbitMQConfigs{
			URL: os.Getenv("RABBITMQ.URL"),
		},
	}

	return config, nil
}
