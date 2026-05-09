package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func TestInsertCartItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Supplier",
		Contact_email: "supplier@testmail.com",
		Contact_phone: "0000000000",
	}
	newSupplier, err = db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProduct := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            1999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
	newProduct, err = db.InsertProduct(newProduct)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
	newAccount, err = db.InsertAccount(newAccount)
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
	newCart, err = db.InsertCart(newCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	newItem := model.Cart_item{
		Cart_id:    newCart.Cart_id,
		Product_id: newProduct.Product_id,
		Quantity:   3,
		Unit_price: 1999,
	}

	item, err := db.InsertCartItem(newItem)

	if err != nil {
		t.Errorf("InsertCartItem() call failed: %s", err.Error())
	}

	if item.Cart_id != newItem.Cart_id {
		t.Errorf(
			"Expected Cart_id %d, Got %d",
			newItem.Cart_id,
			item.Cart_id,
		)
	}

	if item.Product_id != newItem.Product_id {
		t.Errorf(
			"Expected Product_id %d, Got %d",
			newItem.Product_id,
			item.Product_id,
		)
	}
}

func TestGetCartItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Supplier",
		Contact_email: "supplier@testmail.com",
		Contact_phone: "0000000000",
	}
	newSupplier, err = db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProduct := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            1999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
	newProduct, err = db.InsertProduct(newProduct)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
	newAccount, err = db.InsertAccount(newAccount)
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
	newCart, err = db.InsertCart(newCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	insItem := model.Cart_item{
		Cart_id:    newCart.Cart_id,
		Product_id: newProduct.Product_id,
		Quantity:   3,
		Unit_price: 1999,
	}

	newItem, err := db.InsertCartItem(insItem)
	if err != nil {
		t.Fatalf("InsertCartItem() call failed: %s", err.Error())
	}

	getItem, err := db.GetCartItem(newItem.Cart_id, newItem.Product_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getItem.Cart_id != newItem.Cart_id {
		t.Errorf(
			"Got Cart_id %d, Expected %d. Expected GetCartItem to return same Cart_id as InsertCartItem inserted.",
			getItem.Cart_id,
			newItem.Cart_id,
		)
	}

	if getItem.Product_id != newItem.Product_id {
		t.Errorf(
			"Got Product_id %d, Expected %d. Expected GetCartItem to return same Product_id as InsertCartItem inserted.",
			getItem.Product_id,
			newItem.Product_id,
		)
	}
}

func TestGetCartItemsByCart(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Supplier",
		Contact_email: "supplier@testmail.com",
		Contact_phone: "0000000000",
	}
	newSupplier, err = db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProductOne := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Test Product One",
		Description:      "Test Description",
		Price:            1999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
	newProductOne, err = db.InsertProduct(newProductOne)
	if err != nil {
		t.Fatalf("InsertProduct() call failed on first product: %s", err.Error())
	}

	newProductTwo := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Test Product Two",
		Description:      "Test Description",
		Price:            999,
		Discount_percent: 0,
		Stock_quantity:   5,
		Is_active:        true,
	}
	newProductTwo, err = db.InsertProduct(newProductTwo)
	if err != nil {
		t.Fatalf("InsertProduct() call failed on second product: %s", err.Error())
	}

	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
	newAccount, err = db.InsertAccount(newAccount)
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
	newCart, err = db.InsertCart(newCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	_, err = db.InsertCartItem(model.Cart_item{
		Cart_id:    newCart.Cart_id,
		Product_id: newProductOne.Product_id,
		Quantity:   3,
		Unit_price: 1999,
	})
	if err != nil {
		t.Fatalf("InsertCartItem() call failed on first item: %s", err.Error())
	}

	_, err = db.InsertCartItem(model.Cart_item{
		Cart_id:    newCart.Cart_id,
		Product_id: newProductTwo.Product_id,
		Quantity:   1,
		Unit_price: 999,
	})
	if err != nil {
		t.Fatalf("InsertCartItem() call failed on second item: %s", err.Error())
	}

	items, err := db.GetCartItemsByCart(newCart.Cart_id)

	if err != nil {
		t.Errorf("GetCartItemsByCart() call failed: %s", err.Error())
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 cart items, Got %d", len(items))
	}

	for _, item := range items {
		if item.Cart_id != newCart.Cart_id {
			t.Errorf(
				"Expected all items to have Cart_id %d, Got %d",
				newCart.Cart_id,
				item.Cart_id,
			)
		}
	}
}

func TestUpdateCartItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Supplier",
		Contact_email: "supplier@testmail.com",
		Contact_phone: "0000000000",
	}
	newSupplier, err = db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProduct := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            1999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
	newProduct, err = db.InsertProduct(newProduct)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
	newAccount, err = db.InsertAccount(newAccount)
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
	newCart, err = db.InsertCart(newCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	insItem := model.Cart_item{
		Cart_id:    newCart.Cart_id,
		Product_id: newProduct.Product_id,
		Quantity:   3,
		Unit_price: 1999,
	}

	newItem, err := db.InsertCartItem(insItem)
	if err != nil {
		t.Fatalf("InsertCartItem() call failed: %s", err.Error())
	}

	newQuantity := 5
	updatedItem, err := db.UpdateCartItem(newItem.Cart_id, newItem.Product_id, newQuantity)
	if err != nil {
		t.Errorf("UpdateCartItem() call failed: %s", err.Error())
	}

	if updatedItem.Quantity != newQuantity {
		t.Errorf("Expected Quantity %d, Got %d", newQuantity, updatedItem.Quantity)
	}

	if updatedItem.Cart_id != newItem.Cart_id {
		t.Errorf(
			"Expected Cart_id to remain %d after update, Got %d",
			newItem.Cart_id,
			updatedItem.Cart_id,
		)
	}

	if updatedItem.Product_id != newItem.Product_id {
		t.Errorf(
			"Expected Product_id to remain %d after update, Got %d",
			newItem.Product_id,
			updatedItem.Product_id,
		)
	}
}

func TestDeleteCartItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Supplier",
		Contact_email: "supplier@testmail.com",
		Contact_phone: "0000000000",
	}
	newSupplier, err = db.InsertSupplier(newSupplier)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	newProduct := model.Product{
		Supplier_id:      newSupplier.Supplier_id,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            1999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
	newProduct, err = db.InsertProduct(newProduct)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	newAccount := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
	newAccount, err = db.InsertAccount(newAccount)
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
	newCart, err = db.InsertCart(newCart)
	if err != nil {
		t.Fatalf("InsertCart() call failed: %s", err.Error())
	}

	insItem := model.Cart_item{
		Cart_id:    newCart.Cart_id,
		Product_id: newProduct.Product_id,
		Quantity:   3,
		Unit_price: 1999,
	}

	newItem, err := db.InsertCartItem(insItem)
	if err != nil {
		t.Fatalf("InsertCartItem() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteCartItem(newItem.Cart_id, newItem.Product_id)
	if err != nil {
		t.Errorf("DeleteCartItem() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedItem, err := db.GetCartItem(newItem.Cart_id, newItem.Product_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %s", err.Error())
	}

	if deletedItem.Cart_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedItem.Cart_id)
	}
}
