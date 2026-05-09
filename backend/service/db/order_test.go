package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func newTestOrder() (model.Orders, error) {
	newAcc := model.Account{
		Account_id:    -1,
		Full_name:     "Test Name",
		Email:         fmt.Sprintf("%s@testmail.com", time.Now()),
		Password_hash: "somehash",
		Role:          "customer",
		Created_at:    time.Now(),
	}
	acc, err := db.InsertAccount(newAcc)
	if err != nil {
		return model.Orders{}, err
	}

	return model.Orders{
		Order_id:        -1,
		Account_id:      acc.Account_id,
		Order_number:    time.Now().GoString(),
		Status:          "pending",
		Subtotal_amount: 100000,
		Discount_amount: 25000,
		Total_amount:    75000,
		Placed_at:       time.Now(),
	}, nil
}

func TestInsertOrder(t *testing.T) {
	newOrder, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	ord, err := db.InsertOrder(newOrder)

	if err != nil {
		t.Errorf("InsertOrder() call failed: %s", err.Error())
	}

	if ord.Order_id <= 0 {
		t.Errorf("Expected valid Order_id, Got: %d", ord.Order_id)
	}
}

func TestGetOrderById(t *testing.T) {
	insOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := db.InsertOrder(insOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	getOrd, err := db.GetOrderById(newOrd.Order_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getOrd.Order_id != newOrd.Order_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetOrder to return same order id as InsertOrder inserted.",
			getOrd.Order_id,
			newOrd.Order_id,
		)
	}

}

func TestUpdateOrder(t *testing.T) {
	insOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := db.InsertOrder(insOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	chgOrd := newOrd
	chgOrd.Status = "paid"
	rowCount, err := db.UpdateOrder(chgOrd)
	if err != nil {
		t.Errorf("UpdateOrder() call failed: %s", err.Error())
	}

	updatedOrd, err := db.GetOrderById(chgOrd.Order_id)

	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	if updatedOrd.Status != "paid" {
		t.Errorf("Expected 'paid', Got %s", updatedOrd.Status)
	}
}

func TestDeleteOrderById(t *testing.T) {
	insOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := db.InsertOrder(insOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteOrderById(newOrd.Order_id)
	if err != nil {
		t.Errorf("DeleteOrderById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedOrd, err := db.GetOrderById(newOrd.Order_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %s", err.Error())
	}

	if deletedOrd.Order_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedOrd.Order_id)
	}
}
