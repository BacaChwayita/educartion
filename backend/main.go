package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var cfg *config.Config

func handleRegister(w http.ResponseWriter, req *http.Request) {
	cfg.Logs.Logger.Info("Register Request Received", slog.String("remoteAddress", req.RemoteAddr))

	// TODO: Do calls to fulfill Registration
	//       Possibly auth package that returns (status, response)

	// TODO: rpc.WriteJSON(writer, status, response)

}

func handleProducts(w http.ResponseWriter, req *http.Request) {
	cfg.Logs.Logger.Info("Products Request Received", slog.String("remoteAddress", req.RemoteAddr))

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
	http.Handle("/api/auth/register", http.HandlerFunc(handleRegister))
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
