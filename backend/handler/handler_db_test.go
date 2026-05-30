package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupTestDBConfigAndConnection(t *testing.T) {
	t.Helper()

	_ = godotenv.Load("./../.env.test")

	cfg := config.GetConfig()
	cfg.DBConnection.Ctx = context.Background()

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
		cfg.DB.SSLMode,
	)

	var err error
	cfg.DBConnection.Pool, err = pgxpool.New(cfg.DBConnection.Ctx, connectionString)
	if err != nil {
		t.Skipf("Skipping DB-backed handler tests: %s", err.Error())
	}

	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		t.Skipf("Skipping DB-backed handler tests: %s", err.Error())
	}

	t.Cleanup(func() {
		cfg.DBConnection.Pool.Close()
	})
}

func TestHandleRegisterLoginLogoutFlow(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	email := fmt.Sprintf("handler-%d@test.local", time.Now().UnixNano())
	password := "TestPassword123!"

	registerBody, _ := json.Marshal(model.RegisterRequest{
		Full_name:     "Handler Test",
		Email:         email,
		Password_text: password,
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(registerBody))
	HandleRegister(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	loginBody, _ := json.Marshal(model.LoginWithEmailRequest{
		Email:         email,
		Password_text: password,
	})

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	HandleLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var loginResponse model.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	if loginResponse.Account_id <= 0 {
		t.Fatalf("Expected valid account id, got %d", loginResponse.Account_id)
	}
	if loginResponse.Email != email {
		t.Fatalf("Expected email %s, got %s", email, loginResponse.Email)
	}
	if loginResponse.Token == "" {
		t.Fatal("Expected non-empty token")
	}

	logoutBody, _ := json.Marshal(model.LogoutRequest{Token: loginResponse.Token})
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(logoutBody))
	HandleLogout(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestHandleGetOrderByIDInvalidID(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orders/abc", nil)
	req.SetPathValue("id", "abc")

	HandleGetOrderByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
