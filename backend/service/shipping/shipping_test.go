package shipping_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	shippingsvc "github.com/P-SEN371-Group-3/educartion/service/shipping"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupTestDBConfigAndConnection(t *testing.T) {
	_ = godotenv.Load("./../../.env.test")

	cfg := config.GetConfig()
	cfg.DBConnection.Ctx = context.Background()

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
		cfg.DB.SSLMode,
	)

	var err error
	cfg.DBConnection.Pool, err = pgxpool.New(cfg.DBConnection.Ctx, connectionString)
	if err != nil {
		t.Skipf("Skipping DB-backed shipping service tests: %s", err.Error())
	}

	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		t.Skipf("Skipping DB-backed shipping service tests: %s", err.Error())
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

func newTestShipment(orderID int) model.Shipment {
	return model.Shipment{
		Order_id:               orderID,
		Courier_name:           "Test Courier",
		Courier_type:           "ground",
		Delivery_reference:     fmt.Sprintf("REF-%d", time.Now().UnixNano()),
		Delivery_address_line1: "123 Test St",
		Delivery_city:          "Testville",
		Delivery_postal_code:   "12345",
		Delivery_country:       "Testland",
		Recipient_name:         "Jane Doe",
		Shipment_status:        "pending",
		Created_at:             time.Now(),
	}
}

func TestGetShipmentByOrderID(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	account, err := db.InsertAccount(newTestAccount())
	if err != nil {
		t.Fatalf("InsertAccount() failed: %s", err.Error())
	}

	order, err := db.InsertOrder(newTestOrder(account.Account_id))
	if err != nil {
		t.Fatalf("InsertOrder() failed: %s", err.Error())
	}

	shipment, err := db.InsertShipment(newTestShipment(order.Order_id))
	if err != nil {
		t.Fatalf("InsertShipment() failed: %s", err.Error())
	}

	result, err := shippingsvc.GetShipmentByOrderID(order.Order_id)
	if err != nil {
		t.Fatalf("GetShipmentByOrderID() failed: %s", err.Error())
	}

	if result.Shipment_id != shipment.Shipment_id {
		t.Errorf("Expected Shipment_id %d, got %d", shipment.Shipment_id, result.Shipment_id)
	}
	if result.Order_id != order.Order_id {
		t.Errorf("Expected Order_id %d, got %d", order.Order_id, result.Order_id)
	}
}

func TestGetShipmentByOrderIDNotFound(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	_, err := shippingsvc.GetShipmentByOrderID(-1)
	if err == nil {
		t.Fatal("Expected error for missing shipment, got nil")
	}
	if !errors.Is(err, shippingsvc.ErrShipmentNotFound) {
		t.Fatalf("Expected ErrShipmentNotFound, got %v", err)
	}
}
