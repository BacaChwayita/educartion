package db_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// SetupTestDBConfigAndConnection
// A Helper function to set up a db config and connection for tests
// example call: `cfg, err := SetupTestDBConfigAndConnection()
func SetupTestDBConfigAndConnection() (*config.Config, error) {
	//
	// Get DB Connection set up for testing
	//
	var err error
	var cfg config.Config

	_ = godotenv.Load("./../../.env.test")
	cfg = config.LoadConfig()

	cfg.DBConnection.Ctx = context.Background()

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	cfg.DBConnection.Pool, err = pgxpool.New(cfg.DBConnection.Ctx, connectionString)
	if err != nil {
		return &cfg, errors.New("Unable to connect to database: " + err.Error())
	}

	//
	// verify the connection
	//
	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		return &cfg, errors.New("Unable to ping database:" + err.Error())
	}

	return &cfg, nil
}
