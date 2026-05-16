package model

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

type ProductRequest struct{
    Search_name string `json:"search_name"`
	Supplier_id int    `json:"supplier_id"`
	Min_price int    `json:"min_price"`
	Max_price int    `json:"max_price"`
	
}

type ProductResponse struct{
	Products [] Product
    
}

