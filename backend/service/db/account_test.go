package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
)

func TestInsertAccount(t *testing.T) {
	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Password_salt: "somesalt",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	cfg, err := SetupTestDBConfigAndConnection()

	acc, err := db.InsertAccount(cfg, newAccount)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if acc.Account_id == -1 {
		t.Error("Expected valid Account_id, got -1")
	}
}
