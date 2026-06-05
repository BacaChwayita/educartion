package db_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func newTestProduct(supplier_id, category_id int) model.Product {
	return model.Product{
		Product_id:       -1,
		Supplier_id:      supplier_id,
		Category_id:      category_id,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            19999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
}
func TestInsertProduct(t *testing.T) {

	newProduct := model.Product{
		Product_id:       -1,
		Supplier_id:      1,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            19999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf(
			"SetupTestDBConfigAndConnection() call failed: %s",
			err.Error(),
		)
	}

	prod, err := db.InsertProduct(newProduct)

	if err != nil {
		t.Errorf("InsertProduct() call failed: %s", err.Error())
	}

	// Ensure database generated a valid ID
	if prod.Product_id <= 0 {
		t.Errorf(
			"Expected valid Product_id, Got: %d",
			prod.Product_id,
		)
	}
}

func TestGetProductById(t *testing.T) {

	insProd := model.Product{
		Product_id:       -1,
		Supplier_id:      1,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            19999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf(
			"SetupTestDBConfigAndConnection() call failed: %s",
			err.Error(),
		)
	}

	newProd, err := db.InsertProduct(insProd)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	getProd, err := db.GetProductById(newProd.Product_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	// Verify IDs match
	if getProd.Product_id != newProd.Product_id {

		t.Errorf(
			"Got %d, Expected %d. Expected GetProductById() to return same product id.",
			getProd.Product_id,
			newProd.Product_id,
		)
	}
}
func TestUpdateProduct(t *testing.T) {

	insProd := model.Product{
		Product_id:       -1,
		Supplier_id:      1,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            19999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf(
			"SetupTestDBConfigAndConnection() call failed: %s",
			err.Error(),
		)
	}

	newProd, err := db.InsertProduct(insProd)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	// Create updated copy
	chgProd := newProd
	chgProd.Price = 29999
	chgProd.Stock_quantity = 25

	rowCount, err := db.UpdateProduct(chgProd)

	if err != nil {
		t.Errorf("UpdateProduct() call failed: %s", err.Error())
	}

	updatedProd, err := db.GetProductById(chgProd.Product_id)

	// Ensure only one row updated
	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	// Validate updated values
	if updatedProd.Price != 29999 {
		t.Errorf(
			"Expected 29999, Got %d",
			updatedProd.Price,
		)
	}

	if updatedProd.Stock_quantity != 25 {
		t.Errorf(
			"Expected 25, Got %d",
			updatedProd.Stock_quantity,
		)
	}
}
func TestDeleteProductById(t *testing.T) {

	insProd := model.Product{
		Product_id:       -1,
		Supplier_id:      1,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            19999,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf(
			"SetupTestDBConfigAndConnection() call failed: %s",
			err.Error(),
		)
	}

	newProd, err := db.InsertProduct(insProd)
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteProductById(newProd.Product_id)

	if err != nil {
		t.Errorf("DeleteProductById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedProd, err := db.GetProductById(newProd.Product_id)

	// Product should no longer exist
	if err != pgx.ErrNoRows {
		t.Errorf(
			"Expected pgx.ErrNoRows error to be returned, Got: %s",
			err.Error(),
		)
	}

	// Empty struct expected after deletion
	if deletedProd.Product_id != 0 {
		t.Errorf(
			"Expected 0, Got %d",
			deletedProd.Product_id,
		)
	}
}
