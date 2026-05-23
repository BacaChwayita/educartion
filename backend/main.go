package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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
	http.Handle("/api/auth/logout", setHandlerFunc(http.HandlerFunc(handler.HandleLogout)))
	http.Handle("/api/auth/me", setHandlerFunc(http.HandlerFunc(handler.HandleGetUser)))

	// Catalog
	// GET ALL PRODUCTS
	// GET /api/products
	http.Handle(
		"/api/products",
		setHandlerFunc(http.HandlerFunc(handler.HandleGetProducts)),
	)

	// GET PRODUCT BY ID
	// GET /api/products/get?id=1
	http.Handle(
		"/api/products/get",
		setHandlerFunc(http.HandlerFunc(handler.HandleGetProductById)),
	)

	// CREATE PRODUCT
	// POST /api/products/create
	http.Handle(
		"/api/products/create",
		setHandlerFunc(http.HandlerFunc(handler.HandleCreateProduct)),
	)

	// UPDATE PRODUCT
	// PUT /api/products/update?id=1
	http.Handle(
		"/api/products/update",
		setHandlerFunc(http.HandlerFunc(handler.HandleUpdateProduct)),
	)

	// DELETE PRODUCT
	// DELETE /api/products/delete?id=1
	http.Handle(
		"/api/products/delete",
		setHandlerFunc(http.HandlerFunc(handler.HandleDeleteProduct)),
	)

	// Cart
	// TODO: rest of Cart
	http.Handle("/api/cart", setHandlerFunc(http.HandlerFunc(handler.HandleCart)))
	http.Handle("/api/cart/items", setHandlerFunc(http.HandlerFunc(handler.HandleCartItems)))
	http.Handle("/api/cart/item/{productId}", setHandlerFunc(http.HandlerFunc(handler.HandleCartItem)))

	// Orders
	http.Handle("/api/orders/{id}", setHandlerFunc(http.HandlerFunc(handler.HandleGetOrderByID)))
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

		authKey := strings.TrimSpace(r.Header.Get("Authorization"))

		var ctx context.Context
		ctx = context.WithValue(r.Context(), handler.RequestIDKey, reqID)
		ctx = context.WithValue(r.Context(), handler.AuthorisationKey, authKey)

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

//Endpoints for Catalog/ Product 



