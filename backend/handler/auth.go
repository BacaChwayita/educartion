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
	var err error

	var rr model.RegisterRequest
	err = json.NewDecoder(req.Body).Decode(&rr)
	if err != nil {
		logErrDecodeBody("HandleRegister", req, err)
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
		return
	}

	err = auth.Register(rr)
	if err != nil {
		if errors.Is(err, auth.ErrDBInsert) {
			rpc.WriteError(w, 500, "Database error occured while registering user")
			return
		}
		rpc.WriteError(w, 500, "Server error occured while registering user")
		return
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
		logErrDecodeBody("HandleLogin", req, err)
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
	}

	tokenString, acc, err := auth.LoginWithEmail(lr)
	if err != nil && !errors.Is(err, auth.ErrFailedToCreateJWT) {
		reqID, _ := req.Context().Value(RequestIDKey).(string)
		cfg.Logs.Logger.Error(
			"Failed to login newly registered user",
			slog.String("error", err.Error()),
			slog.String("func", "HandleRegister"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("x-request-id", reqID),
		)
		switch {
		case errors.Is(err, auth.ErrAccountInactive):
			rpc.WriteError(w, http.StatusConflict, "Account inactive, must reactivate to login")
			return

		case errors.Is(err, auth.ErrPasswordIncorrect):
			rpc.WriteError(w, http.StatusUnauthorized, "Password incorrect")
			return

		case errors.Is(err, auth.ErrAccountNotFound):
			rpc.WriteError(w, http.StatusNotFound, "Account not found")
			return

		default:
			rpc.WriteError(w, 500, "Server error occured during login process")
			return
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

// handleLogout will logout a user by invalidating the JWT
func HandleLogout(w http.ResponseWriter, req *http.Request) {
	var lr model.LogoutRequest
	err := json.NewDecoder(req.Body).Decode(&lr)
	if err != nil {
		logErrDecodeBody("HandleLogout", req, err)
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
		return
	}

	err = auth.Logout(lr)
	if err == auth.ErrAccountLoginNotFound {
		rpc.WriteError(w, http.StatusNotFound, "No record found in database")
		return
	}

	if err != nil {
		rpc.WriteError(w, 500, "Server error occured during logout process")
		return
	}

	rpc.WriteJSON(w, http.StatusNoContent, nil)
}

// handleGetUser will get return the 'logged in' user details
func HandleGetUser(w http.ResponseWriter, req *http.Request) {
	var gur model.GetUserRequest
	err := json.NewDecoder(req.Body).Decode(&gur)
	if err != nil {
		logErrDecodeBody("HandleGetUser", req, err)
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
		return
	}

}

// ------------
// Helper funcs
// ------------

func logErrDecodeBody(func_name string, req *http.Request, err error) {
	cfg := config.GetConfig()

	reqID, _ := req.Context().Value(RequestIDKey).(string)
	cfg.Logs.Logger.Error(
		"Failed to decode http request body",
		slog.String("error", err.Error()),
		slog.String("func", func_name),
		slog.String("timestamp", time.Now().GoString()),
		slog.String("x-request-id", reqID),
	)
}
