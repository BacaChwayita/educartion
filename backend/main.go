package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool
var ctx = context.Background()
var logger *slog.Logger
var cfg config.Config

func handleRegister(w http.ResponseWriter, req *http.Request) {
	logger.Info("Register Request Received", slog.String("remoteAddress", req.RemoteAddr))

	// TODO: Do calls to fulfill Registration
	//       Possibly auth package that returns (status, response)

	// TODO: rpc.WriteJSON(writer, status, response)

}

func init() {
	_ = godotenv.Load()
	cfg = config.LoadConfig()

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		cfg.Logs.Logger.Error("Unable to connect to database", slog.String("error", err.Error()))
		panic("Unable to connect to database: " + err.Error())
	}

	//
	// verify the connection
	//
	if err := pool.Ping(ctx); err != nil {
		cfg.Logs.Logger.Error("Unable to ping database:", slog.String("error", err.Error()))
		panic("Unable to ping database:" + err.Error())
	}

	cfg.Logs.Logger.Info("Connected to PostgreSQL database!")
}

func main() {

	// TODO: Add config for logging

	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil)) // TODO: Chance nil to handler options - config

	logger.Info("Service Ready")

	// Auth
	http.Handle("/api/auth/register", http.HandlerFunc(handleRegister))

	// Catalog

	// Cart

	// Orders

	// Payments

	// Admin

	http.ListenAndServe(":8080", nil) // TODO: Add port to config

}
