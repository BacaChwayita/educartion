package db_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func newTestCategory() model.Category {
	return model.Category{
		Name:        uuid.New().String(),
		Description: "Test Description",
		Is_active:   true,
	}
}

func TestInsertCategory(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() failed: %s", err.Error())
	}

	newCat := newTestCategory()

	cat, err := db.InsertCategory(newCat)
	if err != nil {
		t.Errorf("InsertCategory() failed: %s", err.Error())
	}

	if cat.Category_id <= 0 {
		t.Errorf("Expected valid Category_id, got %d", cat.Category_id)
	}

	if cat.Name != newCat.Name {
		t.Errorf("Expected Name = Test Category, got %s", cat.Name)
	}
}

func TestGetCategoryById(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() failed: %s", err.Error())
	}

	inserted, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("InsertCategory() failed: %s", err.Error())
	}

	got, err := db.GetCategoryById(inserted.Category_id)
	if err != nil {
		t.Errorf("GetCategoryById() failed: %s", err.Error())
	}

	if got.Category_id != inserted.Category_id {
		t.Errorf("Expected %d, got %d", inserted.Category_id, got.Category_id)
	}

	if got.Name != inserted.Name {
		t.Errorf("Expected %s, got %s", inserted.Name, got.Name)
	}
}

func TestUpdateCategory(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() failed: %s", err.Error())
	}

	inserted, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("InsertCategory() failed: %s", err.Error())
	}

	updated := inserted
	updated.Name = uuid.New().String()
	updated.Description = "Updated Description"
	updated.Is_active = false

	rows, err := db.UpdateCategory(updated)
	if err != nil {
		t.Errorf("UpdateCategory() failed: %s", err.Error())
	}

	if rows != 1 {
		t.Errorf("Expected 1 row affected, got %d", rows)
	}

	got, err := db.GetCategoryById(updated.Category_id)
	if err != nil {
		t.Fatalf("GetCategoryById() failed after update: %s", err.Error())
	}

	if got.Name != updated.Name {
		t.Errorf("Expected Updated Category, got %s", got.Name)
	}

	if got.Is_active != false {
		t.Errorf("Expected Is_active=false, got %v", got.Is_active)
	}
}

func TestDeleteCategoryById(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() failed: %s", err.Error())
	}

	inserted, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("InsertCategory() failed: %s", err.Error())
	}

	rows, err := db.DeleteCategoryById(inserted.Category_id)
	if err != nil {
		t.Errorf("DeleteCategoryById() failed: %s", err.Error())
	}

	if rows != 1 {
		t.Errorf("Expected 1 row affected, got %d", rows)
	}

	_, err = db.GetCategoryById(inserted.Category_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows, got %v", err)
	}
}

func TestGetAllCategories(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() failed: %s", err.Error())
	}

	_, err = db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("InsertCategory() failed: %s", err.Error())
	}

	_, err = db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("InsertCategory() failed: %s", err.Error())
	}

	cats, err := db.GetAllCategories()
	if err != nil {
		t.Errorf("GetAllCategories() failed: %s", err.Error())
	}

	if len(cats) < 2 {
		t.Errorf("Expected at least 2 categories, got %d", len(cats))
	}
}
