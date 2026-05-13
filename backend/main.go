package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/handler"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var cfg *config.Config

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
	http.Handle("/api/auth/register", setHandlerFunc(http.HandlerFunc(handler.HandleRegister)))
	http.Handle("/api/auth/login", setHandlerFunc(http.HandlerFunc(handler.HandleLogin)))
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
// before asking the @h handler function to serve the request
func setHandlerFunc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), handler.RequestIDKey, reqID)

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
		h.ServeHTTP(w, r.WithContext(ctx))
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
