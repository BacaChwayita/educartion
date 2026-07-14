package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/* -----------------------
   helpers
------------------------ */

func req(method, url, body, auth string) *http.Request {
	r := httptest.NewRequest(method, url, strings.NewReader(body))
	if auth != "" {
		r.Header.Set("Authorization", auth)
	}
	return r
}

func TestHandleCart_MethodNotAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	r := req(http.MethodPost, "/api/cart", "", "Bearer token")

	HandleCart(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestGetBearerToken_Missing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := getBearerToken(r)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetBearerToken_InvalidPrefix(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Token abc")

	_, err := getBearerToken(r)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHandleAddCart_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := req(http.MethodPost, "/api/cart/item", "bad-json", "Bearer token")

	HandleAddCart(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleAddCart_InvalidQuantity(t *testing.T) {
	w := httptest.NewRecorder()
	r := req(http.MethodPost, "/api/cart/item", `{"product_id":1,"quantity":0}`, "Bearer token")

	HandleAddCart(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleUpdateCartItem_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	r := req(http.MethodPut, "/api/cart/item/update/abc", `{"quantity":2}`, "Bearer token")

	HandleUpdateCartItem(w, r, "abc")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleDeleteCartItem_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	r := req(http.MethodDelete, "/api/cart/item/delete/-1", "", "Bearer token")

	HandleDeleteCartItem(w, r, "-1")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleCartItem_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/api/cart/item/", nil)

	HandleCartItem(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleCartItem_MethodNotAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/cart/item/1", nil)

	HandleCartItem(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
