package auth_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// SetupTestDBConfigAndConnection
// A Helper function to set up a db config and connection for tests
// This is to replace the fact that the init() function in main.go does not run
func SetupTestDBConfigAndConnection() error {
	//
	// Get DB Connection set up for testing
	//
	var err error
	var cfg *config.Config

	_ = godotenv.Load("./../../.env.test")
	cfg = config.GetConfig()

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
		return errors.New("Unable to connect to database: " + err.Error())
	}

	//
	// verify the connection
	//
	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		return errors.New("Unable to ping database:" + err.Error())
	}

	return nil
}

func TestRegister(t *testing.T) {
	user := model.RegisterRequest{
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_text: "MyPassword123!",
	}

	err := SetupTestDBConfigAndConnection()

	err = auth.Register(user)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}
}

func TestLoginWithEmail(t *testing.T) {
	testEmail := fmt.Sprintf("%s@testmail.com", time.Now())
	testPassword := fmt.Sprintf("%s@testmail.com", time.Now())

	user := model.LoginWithEmailRequest{
		Email:         testEmail,
		Password_text: testPassword,
	}

	registerTestUser := model.RegisterRequest{
		Full_name:     "Test Name",
		Email:         testEmail,
		Password_text: testPassword,
	}

	err := SetupTestDBConfigAndConnection()

	auth.Register(registerTestUser)

	jwt, acc, err := auth.LoginWithEmail(user)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	if jwt == "" {
		t.Errorf("Expected JWT Token, got blank value")
	}

	if acc.Account_id <= 0 {
		t.Errorf("Expected valid account_id, got %d", acc.Account_id)
	}
}
