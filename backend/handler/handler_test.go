package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
)

func decodeErrorMessage(t *testing.T, rr *httptest.ResponseRecorder) string {
	t.Helper()

	var payload map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	return payload["error"]
}

func TestGetBearerToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)

	if _, err := getBearerToken(req); err == nil {
		t.Fatal("Expected error for missing Authorization header, got nil")
	}

	req.Header.Set("Authorization", "Token abc")
	if _, err := getBearerToken(req); err == nil || err.Error() != "authorization header must use Bearer token" {
		t.Fatalf("Expected bearer scheme error, got %v", err)
	}

	req.Header.Set("Authorization", "Bearer ")
	if _, err := getBearerToken(req); err == nil || err.Error() != "authorization header must use Bearer token" {
		t.Fatalf("Expected bearer scheme error, got %v", err)
	}

	req.Header.Set("Authorization", "Bearer test-token")
	token, err := getBearerToken(req)
	if err != nil {
		t.Fatalf("Expected valid token, got %v", err)
	}
	if token != "test-token" {
		t.Fatalf("Expected token test-token, got %s", token)
	}
}

func TestHandleCartMethodNotAllowed(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/cart", nil)

	HandleCart(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleCartItemsMissingAuth(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/cart/items", nil)

	HandleCartItems(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if msg := decodeErrorMessage(t, rr); msg != "authorization header required" {
		t.Fatalf("Expected authorization error, got %q", msg)
	}
}

func TestHandleAddCartInvalidQuantity(t *testing.T) {
	rr := httptest.NewRecorder()
	body, _ := json.Marshal(model.AddCartItemRequest{ProductID: 1, Quantity: 0})
	req := httptest.NewRequest(http.MethodPost, "/api/cart/items", bytes.NewReader(body))

	HandleAddCart(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleUpdateCartItemInvalidProductID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/cart/item/abc", strings.NewReader(`{"quantity":2}`))

	HandleUpdateCartItem(rr, req, "abc")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleUpdateCartItemInvalidQuantity(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/cart/item/1", strings.NewReader(`{"quantity":0}`))

	HandleUpdateCartItem(rr, req, "1")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleDeleteCartItemInvalidProductID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/cart/item/0", nil)

	HandleDeleteCartItem(rr, req, "0")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleCartItemMissingProductID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/cart/item/", nil)

	HandleCartItem(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleCartItemMethodNotAllowed(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/cart/item/1", nil)

	HandleCartItem(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleGetProductsInvalidBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader("{"))

	HandleGetProducts(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if msg := decodeErrorMessage(t, rr); msg != rpc.ErrDecodeHTTPRequestBody {
		t.Fatalf("Expected decode error, got %q", msg)
	}
}

func TestHandleGetProductByIdInvalidQuery(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/product?id=abc", nil)

	HandleGetProductById(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleCreateProductInvalidBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader("{"))

	HandleCreateProduct(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if msg := decodeErrorMessage(t, rr); msg != rpc.ErrDecodeHTTPRequestBody {
		t.Fatalf("Expected decode error, got %q", msg)
	}
}

func TestHandleUpdateProductInvalidID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/products?id=abc", strings.NewReader(`{"name":"bad"}`))

	HandleUpdateProduct(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleDeleteProductInvalidID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/products?id=abc", nil)

	HandleDeleteProduct(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandlePaymentsValidation(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/payments", nil)

	HandlePayments(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/payments", strings.NewReader("{"))
	HandlePayments(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	rr = httptest.NewRecorder()
	body, _ := json.Marshal(map[string]any{
		"order_id":       1,
		"payment_method": "card",
		"amount":         1000,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/payments", bytes.NewReader(body))
	HandlePayments(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestHandlePaymentByIDValidation(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/payments/1", nil)

	HandlePaymentByID(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/payments/abc", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	HandlePaymentByID(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleShipmentByOrderValidation(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/shipments/1", nil)

	HandleShipmentByOrder(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/shipments/abc", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	HandleShipmentByOrder(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleRegisterInvalidBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader("{"))

	HandleRegister(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleLogoutInvalidBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/logout", strings.NewReader("{"))

	HandleLogout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleGetUserInvalidBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader("{"))

	HandleGetUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
