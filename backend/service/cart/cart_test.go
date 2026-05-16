package cart_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/cart"
	"github.com/P-SEN371-Group-3/educartion/service/db"
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
		t.Skipf("Skipping DB-backed cart service tests: %s", err.Error())
	}

	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		t.Skipf("Skipping DB-backed cart service tests: %s", err.Error())
	}
}

func newTestSupplier() model.Supplier {
	return model.Supplier{
		Name:          fmt.Sprintf("Test Supplier %d", time.Now().UnixNano()),
		Contact_email: fmt.Sprintf("supplier-%d@test.local", time.Now().UnixNano()),
		Contact_phone: "555-0100",
	}
}

func newTestProduct(supplierID int) model.Product {
	return model.Product{
		Supplier_id:      supplierID,
		Name:             fmt.Sprintf("Test Product %d", time.Now().UnixNano()),
		Description:      "Service layer test product",
		Price:            1000,
		Discount_percent: 0,
		Stock_quantity:   100,
		Is_active:        true,
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

func TestGetCartByTokenCreatesCart(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	account, err := db.InsertAccount(newTestAccount())
	if err != nil {
		t.Fatalf("InsertAccount() failed: %s", err.Error())
	}

	token := fmt.Sprintf("test-token-%d", time.Now().UnixNano())
	_, err = db.InsertAccountLogin(model.AccountLogin{
		Account_id:   account.Account_id,
		Token_string: token,
		Created_at:   time.Now(),
	})
	if err != nil {
		t.Fatalf("InsertAccountLogin() failed: %s", err.Error())
	}

	cartObj, err := cart.GetCartByToken(token)
	if err != nil {
		t.Fatalf("GetCartByToken() failed: %s", err.Error())
	}

	if cartObj.Account_id != account.Account_id {
		t.Errorf("Expected Account_id %d, got %d", account.Account_id, cartObj.Account_id)
	}

	if cartObj.Session_key != token {
		t.Errorf("Expected Session_key %q, got %q", token, cartObj.Session_key)
	}
}

func TestAddUpdateDeleteCartItemFlow(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	account, err := db.InsertAccount(newTestAccount())
	if err != nil {
		t.Fatalf("InsertAccount() failed: %s", err.Error())
	}

	token := fmt.Sprintf("test-token-%d", time.Now().UnixNano())
	_, err = db.InsertAccountLogin(model.AccountLogin{
		Account_id:   account.Account_id,
		Token_string: token,
		Created_at:   time.Now(),
	})
	if err != nil {
		t.Fatalf("InsertAccountLogin() failed: %s", err.Error())
	}

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() failed: %s", err.Error())
	}

	product, err := db.InsertProduct(newTestProduct(supplier.Supplier_id))
	if err != nil {
		t.Fatalf("InsertProduct() failed: %s", err.Error())
	}

	addReq := model.AddCartItemRequest{
		ProductID: product.Product_id,
		Quantity:  2,
	}

	cartItem, err := cart.AddCartItem(token, addReq)
	if err != nil {
		t.Fatalf("AddCartItem() failed: %s", err.Error())
	}

	if cartItem.Product_id != product.Product_id {
		t.Errorf("Expected Product_id %d, got %d", product.Product_id, cartItem.Product_id)
	}

	if cartItem.Quantity != 2 {
		t.Errorf("Expected Quantity 2, got %d", cartItem.Quantity)
	}

	storedCart, err := db.GetCartById(cartItem.Cart_id)
	if err != nil {
		t.Fatalf("GetCartById() failed: %s", err.Error())
	}

	expectedTotal := product.Price * cartItem.Quantity
	if storedCart.Total_Price != expectedTotal {
		t.Errorf("Expected Total_Price %d, got %d", expectedTotal, storedCart.Total_Price)
	}

	updatedItem, err := cart.UpdateCartItem(token, product.Product_id, 5)
	if err != nil {
		t.Fatalf("UpdateCartItem() failed: %s", err.Error())
	}

	if updatedItem.Quantity != 5 {
		t.Errorf("Expected Quantity 5, got %d", updatedItem.Quantity)
	}

	storedCart, err = db.GetCartById(updatedItem.Cart_id)
	if err != nil {
		t.Fatalf("GetCartById() after update failed: %s", err.Error())
	}

	expectedTotal = product.Price * updatedItem.Quantity
	if storedCart.Total_Price != expectedTotal {
		t.Errorf("Expected Total_Price %d after update, got %d", expectedTotal, storedCart.Total_Price)
	}

	if err := cart.DeleteCartItem(token, product.Product_id); err != nil {
		t.Fatalf("DeleteCartItem() failed: %s", err.Error())
	}

	items, err := cart.GetCartItems(token)
	if err != nil {
		t.Fatalf("GetCartItems() failed: %s", err.Error())
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 cart items after delete, got %d", len(items))
	}

	storedCart, err = db.GetCartById(cartItem.Cart_id)
	if err != nil {
		t.Fatalf("GetCartById() after delete failed: %s", err.Error())
	}

	if storedCart.Total_Price != 0 {
		t.Errorf("Expected Total_Price 0 after deleting item, got %d", storedCart.Total_Price)
	}
}

func TestAddCartItemInvalidQuantity(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	account, err := db.InsertAccount(newTestAccount())
	if err != nil {
		t.Fatalf("InsertAccount() failed: %s", err.Error())
	}

	token := fmt.Sprintf("test-token-%d", time.Now().UnixNano())
	_, err = db.InsertAccountLogin(model.AccountLogin{
		Account_id:   account.Account_id,
		Token_string: token,
		Created_at:   time.Now(),
	})
	if err != nil {
		t.Fatalf("InsertAccountLogin() failed: %s", err.Error())
	}

	_, err = cart.AddCartItem(token, model.AddCartItemRequest{ProductID: 1, Quantity: 0})
	if err == nil {
		t.Fatal("Expected error for zero quantity, got nil")
	}

	if !errors.Is(err, cart.ErrInvalidQuantity) {
		t.Fatalf("Expected ErrInvalidQuantity, got %v", err)
	}
}
