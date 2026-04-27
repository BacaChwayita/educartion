package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	Logs LogConfig
	DB   PostgresConfig
	Port string
}

type LogConfig struct {
	Logger *slog.Logger
	Level  string
}

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func LoadConfig() Config {
	_ = godotenv.Load(".env")

	config := Config{
		Port: os.Getenv("PORT"),
		Logs: LogConfig{
			Logger: nil,
			Level:  os.Getenv("LOG_LEVEL"),
		},
		DB: PostgresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PWD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
	}

	getLogger := func(envChoice string) *slog.Logger {
		handlerOpts := slog.HandlerOptions{}
		switch config.Logs.Level {
		case "DEBUG":
			handlerOpts.Level = slog.LevelDebug
		case "INFO":
			handlerOpts.Level = slog.LevelInfo
		case "WARN":
			handlerOpts.Level = slog.LevelWarn
		case "ERROR":
			handlerOpts.Level = slog.LevelError
		default:
			handlerOpts.Level = slog.LevelInfo
		}

		switch envChoice {
		case "JSON":
			return slog.New(slog.NewJSONHandler(os.Stdout, &handlerOpts))
		case "TEXT":
			return slog.New(slog.NewTextHandler(os.Stdout, &handlerOpts))
		default:
			return slog.New(slog.NewJSONHandler(os.Stdout, &handlerOpts))
		}
	}

	config.Logs.Logger = getLogger(os.Getenv("LOG_TYPE"))

	return config

}
