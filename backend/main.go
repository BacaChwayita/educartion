package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/P-SEN371-Group-3/educartion/config"
)

var logger *slog.Logger
var cfg config.Config

func handleRegister(w http.ResponseWriter, req *http.Request) {
	logger.Info("Register Request Received", slog.String("remoteAddress", req.RemoteAddr))

	// TODO: Do calls to fulfill Registration
	//       Possibly auth package that returns (status, response)

	// TODO: rpc.WriteJSON(writer, status, response)

}

func init() {
	cfg = config.LoadConfig()
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
