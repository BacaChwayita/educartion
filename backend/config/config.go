package config

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Logs         logConfig
	DB           postgresConfig
	DBConnection dbConnection
	Port         string
}

type logConfig struct {
	Logger *slog.Logger
	Level  string
}

type postgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type dbConnection struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

var lock = &sync.Mutex{}
var configInstance *Config

func GetConfig() *Config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil { // Ensure goroutines does not get past previous if while main one is busy locking
			configInstance = loadConfig()
		}
	}

	return configInstance
}

func loadConfig() *Config {

	config := Config{
		Port: os.Getenv("PORT"),
		Logs: logConfig{
			Logger: nil,
			Level:  os.Getenv("LOG_LEVEL"),
		},
		DB: postgresConfig{
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

	return &config

}
