package auth_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/auth"
)

func TestRegister(t *testing.T) {
	user := model.RegisterRequest{
		Full_name:     "Test Name",
		Email:         "test@testmail.com",
		Password_text: "MyPassword123!",
	}

	err := auth.Register(user)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}
}
