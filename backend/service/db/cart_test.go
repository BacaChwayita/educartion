package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func TestInsertCart(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	insAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}

	newAccount, err := db.InsertAccount(insAccount)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	newCart := model.Cart{
		Account_id:  newAccount.Account_id,
		Session_key: fmt.Sprintf("test-session-%s", time.Now()),
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	cart, err := db.InsertCart(newCart)

	if err != nil {
		t.Errorf("InsertCart() call failed: %s", err.Error())
	}

	if cart.Cart_id <= 0 {
		t.Errorf("Expected valid Cart_id, Got: %d", cart.Cart_id)
	}
}

func TestGetCartById(t *testing.T) {
	insCart := model.Cart{
		Account_id:  1,
		Session_key: fmt.Sprintf("test-session-%s", time.Now()),
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newCart, err := db.InsertCart(insCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	getCart, err := db.GetCartById(newCart.Cart_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getCart.Cart_id != newCart.Cart_id {
		t.Errorf(
			"Got Cart_id %d, Expected %d. Expected GetCartById to return same Cart_id as InsertCart inserted.",
			getCart.Cart_id,
			newCart.Cart_id,
		)
	}
}

func TestGetCartByAccountId(t *testing.T) {
	insCart := model.Cart{
		Account_id:  1,
		Session_key: fmt.Sprintf("test-session-%s", time.Now()),
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newCart, err := db.InsertCart(insCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	getCart, err := db.GetCartByAccountId(newCart.Account_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getCart.Account_id != newCart.Account_id {
		t.Errorf(
			"Got Account_id %d, Expected %d. Expected GetCartByAccountId to return same Account_id as InsertCart inserted.",
			getCart.Account_id,
			newCart.Account_id,
		)
	}
}

func TestGetCartBySessionKey(t *testing.T) {
	testSessionKey := fmt.Sprintf("test-session-%s", time.Now())

	insCart := model.Cart{
		Account_id:  1,
		Session_key: testSessionKey,
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newCart, err := db.InsertCart(insCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	getCart, err := db.GetCartBySessionKey(testSessionKey)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getCart.Session_key != newCart.Session_key {
		t.Errorf(
			"Got Session_key %s, Expected %s. Expected GetCartBySessionKey to return same Session_key as InsertCart inserted.",
			getCart.Session_key,
			newCart.Session_key,
		)
	}

	if getCart.Cart_id != newCart.Cart_id {
		t.Errorf(
			"Got Cart_id %d, Expected %d. Expected GetCartBySessionKey to return same Cart_id as InsertCart inserted.",
			getCart.Cart_id,
			newCart.Cart_id,
		)
	}
}

func TestUpdateCart(t *testing.T) {
	insCart := model.Cart{
		Account_id:  1,
		Session_key: fmt.Sprintf("test-session-%s", time.Now()),
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newCart, err := db.InsertCart(insCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	chgCart := newCart
	chgCart.Total_Price = 5997
	chgCart.Updated_at = time.Now()

	rowCount, err := db.UpdateCart(chgCart)
	if err != nil {
		t.Errorf("UpdateCart() call failed: %s", err.Error())
	}

	updatedCart, err := db.GetCartById(chgCart.Cart_id)

	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	if updatedCart.Total_Price != 5997 {
		t.Errorf("Expected Total_Price 5997, Got %d", updatedCart.Total_Price)
	}
}

func TestDeleteCartById(t *testing.T) {
	insCart := model.Cart{
		Account_id:  1,
		Session_key: fmt.Sprintf("test-session-%s", time.Now()),
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newCart, err := db.InsertCart(insCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteCartById(newCart.Cart_id)
	if err != nil {
		t.Errorf("DeleteCartById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedCart, err := db.GetCartById(newCart.Cart_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %s", err.Error())
	}

	if deletedCart.Cart_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedCart.Cart_id)
	}
}

func TestGetCartByJTI(t *testing.T) {
	testToken := fmt.Sprintf("testToken.%s", time.Now())

	insAccount := model.Account{
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

	newAccount, err := db.InsertAccount(insAccount)
	if err != nil {
		t.Fatalf("InsertAccount() call failed: %s", err.Error())
	}

	insCart := model.Cart{
		Account_id:  newAccount.Account_id,
		Session_key: fmt.Sprintf("test-session-%s", time.Now()),
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	newCart, err := db.InsertCart(insCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	insAccountLogin := model.AccountLogin{
		Token_id:     -1,
		Account_id:   newAccount.Account_id,
		Token_string: testToken,
		Created_at:   time.Now(),
	}

	_, err = db.InsertAccountLogin(insAccountLogin)
	if err != nil {
		t.Fatalf("InsertAccountLogin() call failed: %s", err.Error())
	}

	getCart, err := db.GetCartByJTI(testToken)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getCart.Cart_id != newCart.Cart_id {
		t.Errorf(
			"Got Cart_id %d, Expected %d. Expected GetCartByJTI to return same Cart_id as InsertCart inserted.",
			getCart.Cart_id,
			newCart.Cart_id,
		)
	}

	if getCart.Account_id != newAccount.Account_id {
		t.Errorf(
			"Got Account_id %d, Expected %d. Expected GetCartByJTI to return same Account_id as InsertCart inserted.",
			getCart.Account_id,
			newAccount.Account_id,
		)
	}
}
