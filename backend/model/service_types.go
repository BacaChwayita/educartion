package model

type RegisterRequest struct {
	Full_name     string `json:"full_name"`
	Email         string `json:"email"`
	Password_text string `json:"password_hash"`
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
