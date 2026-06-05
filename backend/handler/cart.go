package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	cartsvc "github.com/P-SEN371-Group-3/educartion/service/cart"
	"github.com/jackc/pgx/v5"
)

func HandleCart(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := getBearerToken(req)
	if err != nil {
		rpc.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	cart, err := cartsvc.GetCartByToken(token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			rpc.WriteError(w, http.StatusNotFound, "Cart not found")
			return
		}
		rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving cart")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, cart)
}

func HandleCartItems(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		token, err := getBearerToken(req)
		if err != nil {
			rpc.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		items, err := cartsvc.GetCartItems(token)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				rpc.WriteError(w, http.StatusNotFound, "Cart not found")
				return
			}
			rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving cart items")
			return
		}

		rpc.WriteJSON(w, http.StatusOK, items)
	case http.MethodPost:
		HandleAddCart(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleAddCart(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body model.AddCartItemRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
		return
	}

	if body.ProductID <= 0 || body.Quantity <= 0 {
		rpc.WriteError(w, http.StatusBadRequest, "product_id and quantity must be positive")
		return
	}

	token, err := getBearerToken(req)
	if err != nil {
		rpc.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	cartItem, err := cartsvc.AddCartItem(token, body)
	if err != nil {
		if errors.Is(err, cartsvc.ErrProductNotFound) {
			rpc.WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		if errors.Is(err, cartsvc.ErrInvalidQuantity) {
			rpc.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			rpc.WriteError(w, http.StatusNotFound, "Cart not found")
			return
		}
		rpc.WriteError(w, http.StatusInternalServerError, "Error adding cart item")
		return
	}

	rpc.WriteJSON(w, http.StatusCreated, cartItem)
}

func HandleUpdateCartItem(w http.ResponseWriter, req *http.Request, productID string) {
	if req.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productIDInt, err := strconv.Atoi(productID)
	if err != nil || productIDInt <= 0 {
		rpc.WriteError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	var body model.UpdateCartItemRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		rpc.WriteError(w, http.StatusBadRequest, rpc.ErrDecodeHTTPRequestBody)
		return
	}

	if body.Quantity <= 0 {
		rpc.WriteError(w, http.StatusBadRequest, "quantity must be greater than zero")
		return
	}

	token, err := getBearerToken(req)
	if err != nil {
		rpc.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	updatedItem, err := cartsvc.UpdateCartItem(token, productIDInt, body.Quantity)
	if err != nil {
		if errors.Is(err, cartsvc.ErrCartItemNotFound) {
			rpc.WriteError(w, http.StatusNotFound, "Cart item not found")
			return
		}
		if errors.Is(err, cartsvc.ErrInvalidQuantity) {
			rpc.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			rpc.WriteError(w, http.StatusNotFound, "Cart not found")
			return
		}
		rpc.WriteError(w, http.StatusInternalServerError, "Error updating cart item")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, updatedItem)
}

func HandleDeleteCartItem(w http.ResponseWriter, req *http.Request, productID string) {
	if req.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productIDInt, err := strconv.Atoi(productID)
	if err != nil || productIDInt <= 0 {
		rpc.WriteError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	token, err := getBearerToken(req)
	if err != nil {
		rpc.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cartsvc.DeleteCartItem(token, productIDInt)
	if err != nil {
		if errors.Is(err, cartsvc.ErrCartItemNotFound) {
			rpc.WriteError(w, http.StatusNotFound, "Cart item not found")
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			rpc.WriteError(w, http.StatusNotFound, "Cart not found")
			return
		}
		rpc.WriteError(w, http.StatusInternalServerError, "Error deleting cart item")
		return
	}

	rpc.WriteJSON(w, http.StatusNoContent, nil)
}

func HandleCartItem(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var productID string

	switch {
	case strings.HasPrefix(path, "/api/cart/item/update/"):
		productID = strings.TrimPrefix(path, "/api/cart/item/update/")
	case strings.HasPrefix(path, "/api/cart/item/delete/"):
		productID = strings.TrimPrefix(path, "/api/cart/item/delete/")
	default:
		productID = strings.TrimPrefix(path, "/api/cart/item/")
	}

	if productID == "" {
		http.Error(w, "The product id is missing", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodPut:
		HandleUpdateCartItem(w, r, productID)

	case http.MethodDelete:
		HandleDeleteCartItem(w, r, productID)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getBearerToken(req *http.Request) (string, error) {
	authHeader := strings.TrimSpace(req.Header.Get("Authorization"))
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("authorization header must use Bearer token")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}
