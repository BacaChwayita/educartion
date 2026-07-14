package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupTestDBConfigAndConnection(t *testing.T) {
	t.Helper()

	_ = godotenv.Load("./../.env.test")

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
		t.Skipf("Skipping DB-backed handler tests: %s", err.Error())
	}

	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		t.Skipf("Skipping DB-backed handler tests: %s", err.Error())
	}

	t.Cleanup(func() {
		cfg.DBConnection.Pool.Close()
	})
}

func insertProductImage(t *testing.T, productID int, imageURL string, isPrimary bool) int {
	t.Helper()

	cfg := config.GetConfig()
	var imageID int
	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		`INSERT INTO product_image (product_id, image_url, alt_text, is_primary, sort_order)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING product_image_id`,
		productID,
		imageURL,
		"test image",
		isPrimary,
		1,
	).Scan(&imageID)
	if err != nil {
		t.Fatalf("Failed to insert product image: %s", err.Error())
	}

	return imageID
}

func TestHandleRegisterLoginLogoutFlow(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	email := fmt.Sprintf("handler-%d@test.local", time.Now().UnixNano())
	password := "TestPassword123!"

	registerBody, _ := json.Marshal(model.RegisterRequest{
		Full_name:     "Handler Test",
		Email:         email,
		Password_text: password,
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(registerBody))
	HandleRegister(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	loginBody, _ := json.Marshal(model.LoginWithEmailRequest{
		Email:         email,
		Password_text: password,
	})

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	HandleLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var loginResponse model.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	if loginResponse.Account_id <= 0 {
		t.Fatalf("Expected valid account id, got %d", loginResponse.Account_id)
	}
	if loginResponse.Email != email {
		t.Fatalf("Expected email %s, got %s", email, loginResponse.Email)
	}
	if loginResponse.Token == "" {
		t.Fatal("Expected non-empty token")
	}

	logoutBody, _ := json.Marshal(model.LogoutRequest{Token: loginResponse.Token})
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(logoutBody))
	HandleLogout(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestHandleGetOrderByIDInvalidID(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orders/abc", nil)
	req.SetPathValue("id", "abc")

	HandleGetOrderByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleGetProductsIncludesImageURL(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Handler Supplier",
		Contact_email: "handler@supplier.test",
		Contact_phone: "0000000000",
	}
	newSupplier, err := db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProduct := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Handler Product",
		Description:      "Handler Description",
		Price:            1500,
		Discount_percent: 0,
		Stock_quantity:   5,
		Is_active:        true,
	}
	newProduct, err = db.InsertProduct(newProduct)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	imageURL := "handler-product.jpg"
	insertProductImage(t, newProduct.Product_id, imageURL, true)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/products", bytes.NewReader([]byte(`{}`)))
	HandleGetProducts(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var products []model.Product
	if err := json.NewDecoder(rr.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products response: %v", err)
	}

	var matched *model.Product
	for _, product := range products {
		if product.Product_id == newProduct.Product_id {
			matched = &product
			break
		}
	}

	if matched == nil {
		t.Fatalf("Expected product %d to be in response", newProduct.Product_id)
	}
	if matched.Image_url != imageURL {
		t.Fatalf("Expected image URL %q, got %q", imageURL, matched.Image_url)
	}
}

func TestHandleGetProductByIdIncludesImages(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Handler Supplier Images",
		Contact_email: "handler-images@supplier.test",
		Contact_phone: "0000000001",
	}
	newSupplier, err := db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProduct := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Handler Product Images",
		Description:      "Handler Description",
		Price:            2200,
		Discount_percent: 0,
		Stock_quantity:   3,
		Is_active:        true,
	}
	newProduct, err = db.InsertProduct(newProduct)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	imageURL := "handler-detail.jpg"
	insertProductImage(t, newProduct.Product_id, imageURL, true)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/get?id=%d", newProduct.Product_id), nil)
	HandleGetProductById(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var details model.Product_details
	if err := json.NewDecoder(rr.Body).Decode(&details); err != nil {
		t.Fatalf("Failed to decode product details: %v", err)
	}

	if details.Product_id != newProduct.Product_id {
		t.Fatalf("Expected product id %d, got %d", newProduct.Product_id, details.Product_id)
	}
	if len(details.Product_images) == 0 {
		t.Fatal("Expected product images to be included")
	}
	if details.Product_images[0].Image_url != imageURL {
		t.Fatalf("Expected image URL %q, got %q", imageURL, details.Product_images[0].Image_url)
	}
}
