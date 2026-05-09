package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func TestInsertAccount(t *testing.T) {
	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	acc, err := db.InsertAccount(newAccount)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if acc.Account_id <= 0 {
		t.Errorf("Expected valid Account_id, Got: %d", acc.Account_id)
	}
}

func TestGetAccountById(t *testing.T) {
	insAcc := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newAcc, err := db.InsertAccount(insAcc)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	getAcc, err := db.GetAccountById(newAcc.Account_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getAcc.Account_id != newAcc.Account_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetAccount to return same account id as InsertAccount inserted.",
			getAcc.Account_id,
			newAcc.Account_id,
		)
	}

}

func TestUpdateAccount(t *testing.T) {
	insAcc := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newAcc, err := db.InsertAccount(insAcc)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	chgAcc := newAcc
	chgAcc.Password_hash = "newhash"
	rowCount, err := db.UpdateAccount(chgAcc)
	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	updatedAcc, err := db.GetAccountById(chgAcc.Account_id)

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	if updatedAcc.Password_hash != "newhash" {
		t.Errorf("Expected 'newhash', Got %s", updatedAcc.Password_hash)
	}
}

func TestDeleteAccountById(t *testing.T) {
	insAcc := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newAcc, err := db.InsertAccount(insAcc)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteAccountById(newAcc.Account_id)
	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedAcc, err := db.GetAccountById(newAcc.Account_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %s", err.Error())
	}

	if deletedAcc.Account_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedAcc.Account_id)
	}
}

// TODO:  TestGetAccountByJTI

func TestGetAccountByEmail(t *testing.T) {
	testEmail := fmt.Sprintf("%s@TestGetAccountByEmail.com", time.Now())

	insAcc := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         testEmail,
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newAcc, err := db.InsertAccount(insAcc)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	getAcc, err := db.GetAccountByEmail(testEmail)

	if err != nil {
		t.Errorf("GetAccountByEmail() call failed: %s", err.Error())
	}

	if getAcc.Email != newAcc.Email {
		t.Errorf(
			"Got %d, Expected %d. Expected GetAccount to return same email as InsertAccount inserted.",
			getAcc.Account_id,
			newAcc.Account_id,
		)
	}

	// Theoretically the .Email check above should be enough, but lets double check
	if getAcc.Account_id != newAcc.Account_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetAccount to return same account id as InsertAccount inserted.",
			getAcc.Account_id,
			newAcc.Account_id,
		)
	}

}

func TestIncLoginAttempts(t *testing.T) {
	insAcc := model.Account{
		Account_id:     -1,
		Full_name:      "Test Name",
		Email:          fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash:  "somehash",
		Role:           "customer",
		Login_attempts: 0,
		Is_active:      true,
		Created_at:     time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newAcc, err := db.InsertAccount(insAcc)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	var rowCount int64

	// ========================
	// Increment 1
	// Expect to go from 0 -> 1
	// + active Account
	// ========================

	rowCount, err = db.IncLoginAttempts(newAcc.Account_id)
	if err != nil {
		t.Errorf("IncLoginAttempts() call failed: %s", err.Error())
	}
	if rowCount != 1 {
		t.Errorf("Expected 1 row updated, got %d", rowCount)
	}

	AccIncOnce, err := db.GetAccountById(newAcc.Account_id)
	if err != nil {
		t.Errorf("GetAccountById() call failed: %s", err.Error())
	}

	if AccIncOnce.Login_attempts != 1 {
		t.Errorf("Expected Login_attempts = 1, Got %d", AccIncOnce.Login_attempts)
	}
	if !AccIncOnce.Is_active {
		t.Errorf("Expected Account to remain active, Got Is_active = %t", AccIncOnce.Is_active)
	}

	if t.Failed() {
		t.Fatal("Increment 1 failed, stopping test")
	}

	// ========================
	// Increment 2
	// Expect to go from 1 -> 2
	// + active Account
	// ========================

	rowCount, err = db.IncLoginAttempts(newAcc.Account_id)
	if err != nil {
		t.Errorf("IncLoginAttempts() call failed: %s", err.Error())
	}
	if rowCount != 1 {
		t.Errorf("Expected 1 row updated, got %d", rowCount)
	}

	AccIncOnce, err = db.GetAccountById(newAcc.Account_id)
	if err != nil {
		t.Errorf("GetAccountById() call failed: %s", err.Error())
	}

	if AccIncOnce.Login_attempts != 2 {
		t.Errorf("Expected Login_attempts = 2, Got %d", AccIncOnce.Login_attempts)
	}
	if !AccIncOnce.Is_active {
		t.Errorf("Expected Account to remain active, Got Is_active = %t", AccIncOnce.Is_active)
	}

	if t.Failed() {
		t.Fatal("Increment 2 failed, stopping test")
	}

	// ========================
	// Increment 3
	// Expect to go from 2 -> 3
	// + inactive Account
	// ========================

	rowCount, err = db.IncLoginAttempts(newAcc.Account_id)
	if err != nil {
		t.Errorf("IncLoginAttempts() call failed: %s", err.Error())
	}
	if rowCount != 1 {
		t.Errorf("Expected 1 row updated, got %d", rowCount)
	}

	AccIncOnce, err = db.GetAccountById(newAcc.Account_id)
	if err != nil {
		t.Errorf("GetAccountById() call failed: %s", err.Error())
	}

	if AccIncOnce.Login_attempts != 3 {
		t.Errorf("Expected Login_attempts = 3, Got %d", AccIncOnce.Login_attempts)
	}
	if AccIncOnce.Is_active {
		t.Errorf("Expected Account to deactivate, Got Is_active = %t", AccIncOnce.Is_active)
	}

	if t.Failed() {
		t.Fatal("Increment 3 failed, stopping test")
	}
}

func TestInsertAccountLogin(t *testing.T) {
	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	acc, err := db.InsertAccount(newAccount)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	testToken := fmt.Sprintf("testToken.%s", time.Now())

	newAccountLogin := model.AccountLogin{
		Token_id:     -1,
		Account_id:   acc.Account_id,
		Token_string: testToken,
		Created_at:   time.Now(),
	}

	accLogin, err := db.InsertAccountLogin(newAccountLogin)
	if err != nil {
		t.Errorf("InsertAccountLogin() call failed: %s", err.Error())
	}

	if accLogin.Token_id <= 0 {
		t.Errorf("Expected valid token_id, Got %d", accLogin.Token_id)
	}

}
