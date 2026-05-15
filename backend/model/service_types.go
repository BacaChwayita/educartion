package model

import "time"

type RegisterRequest struct {
	Full_name     string `json:"full_name"`
	Email         string `json:"email"`
	Password_text string `json:"password_hash"`
}

type RegisterResponse struct {
	Account_id int    `json:"account_id"`
	Full_name  string `json:"full_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Token      string `json:"token"`
}

type LoginWithEmailRequest struct {
	Email         string `json:"email"`
	Password_text string `json:"password_hash"`
}

type LoginWithTokenRequest struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	Account_id int    `json:"account_id"`
	Full_name  string `json:"full_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Token      string `json:"token"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type GetUserRequest struct {
	Token string `json:"token"`
}

type GetOrderByIDRequest struct {
	Id int `json:"id"`
}

type GetOrderByIDResponse struct {
	Order_id        int       `json:"order_id"`
	Order_number    string    `json:"order_number"`
	Status          string    `json:"status"`
	Subtotal_amount int       `json:"subtotal_amount"`
	Discount_amount int       `json:"discount_amount"`
	Total_amount    int       `json:"total_amount"`
	Placed_at       time.Time `json:"placed_at"`

	Order_item []Order_item_list `json:"order_item"`
}

type Order_item_list struct {
	Quantity        int `json:"quantity"`
	Unit_price      int `json:"unit_price"`
	Discount_amount int `json:"discount_amount"`

	Product Product_details `json:"product_details"`
}

type Product_details struct {
	Product_id       int    `json:"product_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Price            int    `json:"price"`
	Discount_percent int    `json:"discount_percent"`
	Stock_quantity   int    `json:"stock_quantity"`
	Is_active        bool   `json:"is_active"`

	Supplier Supplier_details `json:"supplier_details"`

	// TODO: Product Image part - need mime/multipart stuff, but for now, ignore
}

type Supplier_details struct {
	Supplier_id int    `json:"supplier_id"`
	Name        string `json:"name"`
}
