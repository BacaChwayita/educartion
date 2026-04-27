package config

import "os"

type Config struct {
	Port string
}

func LoadConfig() Config {
	return Config{
		Port: os.Getenv("PORT"),
	}
}
