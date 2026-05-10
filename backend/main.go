package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	"github.com/P-SEN371-Group-3/educartion/server"
	"github.com/P-SEN371-Group-3/educartion/service/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var cfg *config.Config

// handleRegister will register a user
func handleRegister(w http.ResponseWriter, req *http.Request) {
	var err error

	var rr model.RegisterRequest
	err = json.NewDecoder(req.Body).Decode(&rr)
	if err != nil {
		cfg.Logs.Logger.Error(
			"Failed to decode http request body",
			slog.String("error", err.Error()),
			slog.String("func", "handleRegister"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("x-request-id", req.Context().Value(server.RequestIDKey).(string)),
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
func handleLogin(w http.ResponseWriter, req *http.Request) {

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
			slog.String("x-request-id", req.Context().Value(server.RequestIDKey).(string)),
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
			slog.String("x-request-id", req.Context().Value(server.RequestIDKey).(string)),
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

func handleProducts(w http.ResponseWriter, req *http.Request) {

	// TODO: Products processing

	// TODO: rpc.WriteJSON(writer, status, response)

}

func init() {
	var err error

	_ = godotenv.Load()
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
		cfg.Logs.Logger.Error("Unable to connect to database", slog.String("error", err.Error()))
		panic("Unable to connect to database: " + err.Error())
	}

	//
	// verify the connection
	//
	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		cfg.Logs.Logger.Error("Unable to ping database:", slog.String("error", err.Error()))
		panic("Unable to ping database:" + err.Error())
	}

	cfg.Logs.Logger.Info("Connected to PostgreSQL database!")
}

func main() {

	cfg.Logs.Logger.Info("Service Ready")

	// Auth
	http.Handle("/api/auth/register", setHandlerFunc(http.HandlerFunc(handleRegister)))
	http.Handle("/api/auth/login", setHandlerFunc(http.HandlerFunc(handleLogin)))
	// TODO: rest of Auth

	// Catalog
	http.Handle("/api/products", http.HandlerFunc(handleProducts))
	// TODO: rest of Catalog

	// Cart
	// TODO: rest of Cart

	// Orders
	// TODO: rest of Orders

	// Payments
	// TODO: rest of Payments

	// Admin
	// TODO: rest of Admin

	portString := fmt.Sprintf(":%s", cfg.Port)
	http.ListenAndServe(portString, nil)

}

// setHandlerFunc is a wrapper function that will set context and log requests
// before asking the @handler function to serve the request
func setHandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), server.RequestIDKey, reqID)

		start := time.Now()

		// Logging before the request has been processed
		cfg.Logs.Logger.Info(
			"Request received",
			slog.String("func", "setHandlerFunc"),
			slog.String("timestamp", start.GoString()),
			slog.String("x-request-id", reqID),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
			slog.String("remoteAddress", r.RemoteAddr),
		)

		// Call the actual handler and time it
		handler.ServeHTTP(w, r.WithContext(ctx))
		end := time.Now()

		// Logging after request has been processed
		cfg.Logs.Logger.Info(
			"Request completed",
			slog.String("func", "setHandlerFunc"),
			slog.String("timestamp", end.GoString()),
			slog.Duration("duration", end.Sub(start)),
			slog.String("x-request-id", reqID),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
			slog.String("remoteAddress", r.RemoteAddr),
		)
	})
}
