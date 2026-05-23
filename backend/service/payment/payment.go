package payment

import (
	"errors"
	"strings"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RecordPaymentRequest struct {
	OrderID       int    `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        int    `json:"amount"`
}

func RecordPayment(req RecordPaymentRequest) (model.Payment, error) {
	if req.OrderID <= 0 || req.Amount <= 0 || strings.TrimSpace(req.PaymentMethod) == "" {
		return model.Payment{}, ErrInvalidPaymentRequest
	}

	_, err := db.GetOrderById(req.OrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Payment{}, ErrOrderNotFound
		}
		return model.Payment{}, err
	}

	_, err = db.GetPaymentByOrderId(req.OrderID)
	if err == nil {
		return model.Payment{}, ErrPaymentAlreadyExists
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.Payment{}, err
	}

	return db.InsertPayment(model.Payment{
		Order_id:        req.OrderID,
		Payment_method:  strings.TrimSpace(req.PaymentMethod),
		Payment_status:  "pending",
		Transaction_ref: uuid.NewString(),
		Amount:          req.Amount,
		Paid_at:         time.Now(),
	})
}

func GetPaymentByID(paymentID int) (model.Payment, error) {
	payment, err := db.GetPaymentById(paymentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Payment{}, ErrPaymentNotFound
		}
		return model.Payment{}, err
	}
	return payment, nil
}

func GetPaymentByOrderID(orderID int) (model.Payment, error) {
	payment, err := db.GetPaymentByOrderId(orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Payment{}, ErrPaymentNotFound
		}
		return model.Payment{}, err
	}
	return payment, nil
}
