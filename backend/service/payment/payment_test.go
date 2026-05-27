package payment_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	paymentsvc "github.com/P-SEN371-Group-3/educartion/service/payment"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupTestDBConfigAndConnection(t *testing.T) {
	_ = godotenv.Load("./../../.env.test")

	cfg := config.GetConfig()
	cfg.DBConnection.Ctx = context.Background()

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	var err error
	cfg.DBConnection.Pool, err = pgxpool.New(cfg.DBConnection.Ctx, connectionString)
	if err != nil {
		t.Skipf("Skipping DB-backed payment service tests: %s", err.Error())
	}

	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		t.Skipf("Skipping DB-backed payment service tests: %s", err.Error())
	}
}

func newTestAccount() model.Account {
	return model.Account{
		Full_name:     fmt.Sprintf("Test Account %d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("account-%d@test.local", time.Now().UnixNano()),
		Password_hash: "test-password-hash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
}

func newTestOrder(accountID int) model.Orders {
	return model.Orders{
		Account_id:      accountID,
		Order_number:    fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		Status:          "pending",
		Subtotal_amount: 1000,
		Discount_amount: 0,
		Total_amount:    1000,
		Placed_at:       time.Now(),
	}
}

func TestRecordPayment(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	account, err := db.InsertAccount(newTestAccount())
	if err != nil {
		t.Fatalf("InsertAccount() failed: %s", err.Error())
	}

	order, err := db.InsertOrder(newTestOrder(account.Account_id))
	if err != nil {
		t.Fatalf("InsertOrder() failed: %s", err.Error())
	}

	payment, err := paymentsvc.RecordPayment(paymentsvc.RecordPaymentRequest{
		OrderID:       order.Order_id,
		PaymentMethod: "card",
		Amount:        1000,
	})
	if err != nil {
		t.Fatalf("RecordPayment() failed: %s", err.Error())
	}

	if payment.Order_id != order.Order_id {
		t.Errorf("Expected Order_id %d, got %d", order.Order_id, payment.Order_id)
	}
	if payment.Payment_method != "card" {
		t.Errorf("Expected Payment_method card, got %s", payment.Payment_method)
	}
	if payment.Payment_status != "pending" {
		t.Errorf("Expected Payment_status pending, got %s", payment.Payment_status)
	}
}

func TestRecordPaymentInvalidAmount(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	_, err := paymentsvc.RecordPayment(paymentsvc.RecordPaymentRequest{
		OrderID:       1,
		PaymentMethod: "card",
		Amount:        0,
	})
	if err == nil {
		t.Fatal("Expected error for invalid amount, got nil")
	}
	if !errors.Is(err, paymentsvc.ErrInvalidPaymentRequest) {
		t.Fatalf("Expected ErrInvalidPaymentRequest, got %v", err)
	}
}

func TestGetPaymentByIDNotFound(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	_, err := paymentsvc.GetPaymentByID(-1)
	if err == nil {
		t.Fatal("Expected error for missing payment, got nil")
	}
	if !errors.Is(err, paymentsvc.ErrPaymentNotFound) {
		t.Fatalf("Expected ErrPaymentNotFound, got %v", err)
	}
}

func TestGetPaymentByOrderID(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	account, err := db.InsertAccount(newTestAccount())
	if err != nil {
		t.Fatalf("InsertAccount() failed: %s", err.Error())
	}

	order, err := db.InsertOrder(newTestOrder(account.Account_id))
	if err != nil {
		t.Fatalf("InsertOrder() failed: %s", err.Error())
	}

	recorded, err := paymentsvc.RecordPayment(paymentsvc.RecordPaymentRequest{
		OrderID:       order.Order_id,
		PaymentMethod: "card",
		Amount:        1000,
	})
	if err != nil {
		t.Fatalf("RecordPayment() failed: %s", err.Error())
	}

	loaded, err := paymentsvc.GetPaymentByOrderID(order.Order_id)
	if err != nil {
		t.Fatalf("GetPaymentByOrderID() failed: %s", err.Error())
	}

	if loaded.Payment_id != recorded.Payment_id {
		t.Fatalf("Expected Payment_id %d, got %d", recorded.Payment_id, loaded.Payment_id)
	}
}
