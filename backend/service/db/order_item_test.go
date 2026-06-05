package db_test

import (
	"errors"
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func newTestOrderItem(order_id, product_id int) model.Order_item {
	return model.Order_item{
		Order_id:        order_id,
		Product_id:      product_id,
		Quantity:        2,
		Unit_price:      10,
		Discount_amount: 2,
	}
}

func setupSupplierAndProduct() (model.Supplier, model.Product, error) {

	sup, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		return model.Supplier{}, model.Product{}, errors.New("InsertSupplier() call failed: " + err.Error())
	}

	cat, err := db.InsertCategory(newTestCategory())
	if err != nil {
		return model.Supplier{}, model.Product{}, errors.New("InsertCategory() call failed: " + err.Error())
	}

	prd, err := db.InsertProduct(newTestProduct(sup.Supplier_id, cat.Category_id))
	if err != nil {
		return model.Supplier{}, model.Product{}, errors.New("InsertProduct() call failed: " + err.Error())
	}

	return sup, prd, nil
}

func TestInsertOrderItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	parentOrd, err := db.InsertOrder(newOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	_, prd, err := setupSupplierAndProduct()
	if err != nil {
		t.Fatal(err.Error())
	}

	oi, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id, prd.Product_id))
	if err != nil {
		t.Errorf("InsertOrderItem() call failed: %s", err.Error())
	}

	if oi.Order_id <= 0 {
		t.Errorf("Expected valid Order_id, Got: %d", oi.Order_id)
	}

	if oi.Product_id <= 0 {
		t.Errorf("Expected valid Product_id, Got: %d", oi.Product_id)
	}
}

func TestGetOrderItemById(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	parentOrd, err := db.InsertOrder(newOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	_, prd, err := setupSupplierAndProduct()
	if err != nil {
		t.Fatal(err.Error())
	}

	newOI, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id, prd.Product_id))
	if err != nil {
		t.Fatalf("InsertOrderItem() call failed: %s", err.Error())
	}

	getOI, err := db.GetOrderItemById(newOI.Order_id, newOI.Product_id)
	if err != nil {
		t.Errorf("GetOrderItemById() call failed: %s", err.Error())
	}

	if getOI.Order_id != newOI.Order_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetOrderItemById to return same id as InsertOrderItem inserted.",
			getOI.Order_id,
			newOI.Order_id,
		)
	}
	if getOI.Product_id != newOI.Product_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetOrderItemById to return same id as InsertOrderItem inserted.",
			getOI.Product_id,
			newOI.Product_id,
		)
	}
}

func TestUpdateOrderItem(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	parentOrd, err := db.InsertOrder(newOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	_, prd, err := setupSupplierAndProduct()
	if err != nil {
		t.Fatal(err.Error())
	}

	newOI, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id, prd.Product_id))
	if err != nil {
		t.Fatalf("InsertOrderItem() call failed: %s", err.Error())
	}

	chgOI := newOI
	chgOI.Quantity = 10
	chgOI.Unit_price = 1999

	rowCount, err := db.UpdateOrderItem(chgOI)
	if err != nil {
		t.Errorf("UpdateOrderItem() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	updatedOI, err := db.GetOrderItemById(chgOI.Order_id, chgOI.Product_id)
	if err != nil {
		t.Errorf("GetOrderItemById() after update failed: %s", err.Error())
	}

	if updatedOI.Quantity != 10 {
		t.Errorf("Expected Quantity = 10, Got %d", updatedOI.Quantity)
	}

	if updatedOI.Unit_price != 1999 {
		t.Errorf("Expected Unit_price = 19.99, Got %d", updatedOI.Unit_price)
	}
}

func TestDeleteOrderItemById(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newOrd, err := newTestOrder()
	if err != nil {
		t.Fatal(err.Error())
	}

	parentOrd, err := db.InsertOrder(newOrd)
	if err != nil {
		t.Fatalf("InsertOrder() call failed: %s", err.Error())
	}

	_, prd, err := setupSupplierAndProduct()
	if err != nil {
		t.Fatal(err.Error())
	}

	newOI, err := db.InsertOrderItem(newTestOrderItem(parentOrd.Order_id, prd.Product_id))
	if err != nil {
		t.Fatalf("InsertOrderItem() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteOrderItemById(newOI.Order_id, newOI.Product_id)
	if err != nil {
		t.Errorf("DeleteOrderItemById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedOI, err := db.GetOrderItemById(newOI.Order_id, newOI.Product_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %v", err)
	}

	if deletedOI.Order_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedOI.Order_id)
	}

	if deletedOI.Product_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedOI.Product_id)
	}
}
