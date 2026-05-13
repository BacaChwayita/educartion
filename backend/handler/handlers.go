package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	"github.com/P-SEN371-Group-3/educartion/service/auth"
)

// handleRegister will register a user
func HandleRegister(w http.ResponseWriter, req *http.Request) {
	cfg := config.GetConfig()
	var err error

	var rr model.RegisterRequest
	err = json.NewDecoder(req.Body).Decode(&rr)
	if err != nil {
		cfg.Logs.Logger.Error(
			"Failed to decode http request body",
			slog.String("error", err.Error()),
			slog.String("func", "handleRegister"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("x-request-id", req.Context().Value(RequestIDKey).(string)),
		)
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
	}

	err = auth.Register(rr)
	if errors.Is(err, auth.ErrDBInsert) {
		rpc.WriteError(w, 500, "Database error occured while registering user")
	} else if err != nil {
		rpc.WriteError(w, 500, "Server error occured while registering user")
	}

	rpc.WriteJSON(w, http.StatusCreated, nil)
}

// handleLogin will login a user either via email or through a JWT
func HandleLogin(w http.ResponseWriter, req *http.Request) {
	cfg := config.GetConfig()

	// TODO: Add the JWT login and refactor the below to only be called if not JWT login.

	var err error

	var lr model.LoginWithEmailRequest
	err = json.NewDecoder(req.Body).Decode(&lr)
	if err != nil {
		cfg.Logs.Logger.Error(
			"Failed to decode http request body",
			slog.String("error", err.Error()),
			slog.String("func", "handleLogin"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("x-request-id", req.Context().Value(RequestIDKey).(string)),
		)
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
	}

	tokenString, acc, err := auth.LoginWithEmail(lr)
	if err != nil && !errors.Is(err, auth.ErrFailedToCreateJWT) {
		cfg.Logs.Logger.Error(
			"Failed to login newly registered user",
			slog.String("error", err.Error()),
			slog.String("func", "handleRegister"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("x-request-id", req.Context().Value(RequestIDKey).(string)),
		)
		switch {
		case errors.Is(err, auth.ErrAccountInactive):
			rpc.WriteError(w, http.StatusConflict, "Account inactive, must reactivate to login")

		case errors.Is(err, auth.ErrPasswordIncorrect):
			rpc.WriteError(w, http.StatusUnauthorized, "Password incorrect")

		case errors.Is(err, auth.ErrAccountNotFound):
			rpc.WriteError(w, http.StatusNotFound, "Account not found")

		default:
			rpc.WriteError(w, 500, "Server error occured during login process")
		}
	}

	response := model.LoginResponse{
		Account_id: acc.Account_id,
		Full_name:  acc.Full_name,
		Email:      acc.Email,
		Role:       acc.Role,
		Token:      tokenString,
	}

	rpc.WriteJSON(w, http.StatusOK, response)

}
