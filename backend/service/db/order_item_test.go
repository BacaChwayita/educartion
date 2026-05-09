package db_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func newTestOrderItem(order_id int) model.OrderItem {
	return model.OrderItem{
		Order_item_id: -1,
		Order_id:      order_id,
		Product_name:  "Test Product",
		Quantity:      2,
		Unit_price:    9.99,
	}
}

func TestInsertOrderItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	parentOrd, err := db.InsertOrder(newTestOrder())
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	oi, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id))
	if err != nil {
		t.Errorf("InsertOrderItem() call failed: %s", err.Error())
	}

	if oi.Order_item_id <= 0 {
		t.Errorf("Expected valid Order_item_id, Got: %d", oi.Order_item_id)
	}
}

func TestGetOrderItemById(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	parentOrd, err := db.InsertOrder(newTestOrder())
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	newOI, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id))
	if err != nil {
		t.Fatalf("InsertOrderItem() call failed: %s", err.Error())
	}

	getOI, err := db.GetOrderItemById(newOI.Order_item_id)
	if err != nil {
		t.Errorf("GetOrderItemById() call failed: %s", err.Error())
	}

	if getOI.Order_item_id != newOI.Order_item_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetOrderItemById to return same id as InsertOrderItem inserted.",
			getOI.Order_item_id,
			newOI.Order_item_id,
		)
	}
}

func TestUpdateOrderItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	parentOrd, err := db.InsertOrder(newTestOrder())
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	newOI, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id))
	if err != nil {
		t.Fatalf("InsertOrderItem() call failed: %s", err.Error())
	}

	chgOI := newOI
	chgOI.Quantity = 10
	chgOI.Unit_price = 19.99

	rowCount, err := db.UpdateOrderItem(chgOI)
	if err != nil {
		t.Errorf("UpdateOrderItem() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	updatedOI, err := db.GetOrderItemById(chgOI.Order_item_id)
	if err != nil {
		t.Errorf("GetOrderItemById() after update failed: %s", err.Error())
	}

	if updatedOI.Quantity != 10 {
		t.Errorf("Expected Quantity = 10, Got %d", updatedOI.Quantity)
	}

	if updatedOI.Unit_price != 19.99 {
		t.Errorf("Expected Unit_price = 19.99, Got %f", updatedOI.Unit_price)
	}
}

func TestDeleteOrderItemById(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	parentOrd, err := db.InsertOrder(newTestOrder())
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	newOI, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id))
	if err != nil {
		t.Fatalf("InsertOrderItem() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteOrderItemById(newOI.Order_item_id)
	if err != nil {
		t.Errorf("DeleteOrderItemById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedOI, err := db.GetOrderItemById(newOI.Order_item_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %v", err)
	}

	if deletedOI.Order_item_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedOI.Order_item_id)
	}
}
