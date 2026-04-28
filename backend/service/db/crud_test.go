package db_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func SetupTestDBConfigAndConnection() (*config.Config, error) {
	//
	// Get DB Connection set up for testing
	//
	var err error
	var cfg config.Config

	_ = godotenv.Load(".env.test")
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

func TestInsertAccount(t *testing.T) {
	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         "test@testmail.com",
		Password_hash: "somehash",
		Password_salt: "somesalt",
		Role:          "USER",
		Created_at:    time.Now(),
	}

	cfg, err := SetupTestDBConfigAndConnection()

	acc, err := db.InsertAccount(cfg, newAccount)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if acc.Account_id == -1 {
		t.Error("Expected valid Account_id, got -1")
	}
}
