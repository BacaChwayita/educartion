package model

type RegisterRequest struct {
	Full_name     string `json:"full_name"`
	Email         string `json:"email"`
	Password_text string `json:"password_hash"`
}
