package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreatePaymentRequest struct {
	OrderID       int    `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        int    `json:"amount"`
}

func HandlePayments(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body CreatePaymentRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
		return
	}

	if body.OrderID <= 0 || body.Amount <= 0 || strings.TrimSpace(body.PaymentMethod) == "" {
		rpc.WriteError(w, http.StatusBadRequest, "order_id, payment_method, and amount are required")
		return
	}

	_, err := getBearerToken(req)
	if err != nil {
		rpc.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	payment, err := db.InsertPayment(model.Payment{
		Order_id:        body.OrderID,
		Payment_method:  strings.TrimSpace(body.PaymentMethod),
		Payment_status:  "pending",
		Transaction_ref: uuid.NewString(),
		Amount:          body.Amount,
		Paid_at:         time.Now(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			rpc.WriteError(w, http.StatusNotFound, "Order not found")
			return
		}
		rpc.WriteError(w, http.StatusInternalServerError, "Error recording payment")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, payment)
}

func HandlePaymentByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := getBearerToken(req)
	if err != nil {
		rpc.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	paymentID := strings.TrimPrefix(req.URL.Path, "/api/payments/")
	if paymentID == "" {
		rpc.WriteError(w, http.StatusBadRequest, "payment id is required")
		return
	}

	id, err := strconv.Atoi(paymentID)
	if err != nil || id <= 0 {
		rpc.WriteError(w, http.StatusBadRequest, "invalid payment id")
		return
	}

	payment, err := db.GetPaymentById(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			rpc.WriteError(w, http.StatusNotFound, "Payment not found")
			return
		}
		rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving payment")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, payment)
}
