package db_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func newTestSupplier() model.Supplier {
	return model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Name",
		Contact_email: "test@testmail.com",
		Contact_phone: "somehash",
	}
}
func TestInsertSupplier(t *testing.T) {
	newSupplier := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Name",
		Contact_email: "test@testmail.com",
		Contact_phone: "somehash",
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	sup, err := db.InsertSupplier(newSupplier)

	if err != nil {
		t.Errorf("InsertSupplier() call failed: %s", err.Error())
	}

	if sup.Supplier_id <= 0 {
		t.Errorf("Expected valid Supplier_id, Got: %d", sup.Supplier_id)
	}
}

func TestGetSuplierById(t *testing.T) {
	insSup := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Name",
		Contact_email: "test@testmail.com",
		Contact_phone: "somehash",
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSup, err := db.InsertSupplier(insSup)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	getSup, err := db.GetSupplierById(newSup.Supplier_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getSup.Supplier_id != newSup.Supplier_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetSupplier to return same supplier id as InsertSupplier inserted.",
			getSup.Supplier_id,
			newSup.Supplier_id,
		)
	}

}

func TestUpdateSupplier(t *testing.T) {
	insSup := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Name",
		Contact_email: "test@testmail.com",
		Contact_phone: "somehash",
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSup, err := db.InsertSupplier(insSup)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	chgSup := newSup
	chgSup.Contact_email = "updated@testmail.com"
	rowCount, err := db.UpdateSupplier(chgSup)
	if err != nil {
		t.Errorf("UpdateSupplier() call failed: %s", err.Error())
	}

	updatedSup, err := db.GetSupplierById(chgSup.Supplier_id)

	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	if updatedSup.Contact_email != "updated@testmail.com" {
		t.Errorf("Expected 'updated@testmail.com', Got %s", updatedSup.Contact_email)
	}
}

func TestDeleteSupplierById(t *testing.T) {
	insSup := model.Supplier{
		Supplier_id:   -1,
		Name:          "Test Name",
		Contact_email: "test@testmail.com",
		Contact_phone: "somehash",
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newSup, err := db.InsertSupplier(insSup)
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteSupplierById(newSup.Supplier_id)
	if err != nil {
		t.Errorf("DeleteSupplierById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedSup, err := db.GetSupplierById(newSup.Supplier_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %s", err.Error())
	}

	if deletedSup.Supplier_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedSup.Supplier_id)
	}
}
