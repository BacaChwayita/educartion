package catalog_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/catalog"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupTestDBConfigAndConnection(t *testing.T) {
	_ = godotenv.Load("./../../.env.test")

	cfg := config.GetConfig()
	cfg.DBConnection.Ctx = context.Background()

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	var err error
	cfg.DBConnection.Pool, err = pgxpool.New(cfg.DBConnection.Ctx, connectionString)
	if err != nil {
		t.Skipf("Skipping DB-backed catalog service tests: %s", err.Error())
	}

	if err = cfg.DBConnection.Pool.Ping(cfg.DBConnection.Ctx); err != nil {
		t.Skipf("Skipping DB-backed catalog service tests: %s", err.Error())
	}
}

func newTestSupplier() model.Supplier {
	return model.Supplier{
		Name:          fmt.Sprintf("Test Supplier %d", time.Now().UnixNano()),
		Contact_email: fmt.Sprintf("supplier-%d@test.local", time.Now().UnixNano()),
		Contact_phone: "555-0100",
	}
}

func newTestCategory() model.Category {
	return model.Category{
		Name:        uuid.New().String(),
		Description: "Test Description",
		Is_active:   true,
	}
}

func newTestProduct(supplier_id, category_id int) model.Product {
	return model.Product{
		Product_id:       -1,
		Supplier_id:      supplier_id,
		Category_id:      category_id,
		Name:             "Test Product",
		Description:      "Test Description",
		Price:            1000,
		Discount_percent: 0,
		Stock_quantity:   10,
		Is_active:        true,
	}
}

func TestInsertAndGetProductById(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() failed: %s", err.Error())
	}

	cat, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("newTestCategory() failed: %s", err.Error())
	}

	product := newTestProduct(supplier.Supplier_id, cat.Category_id)
	created, err := catalog.InsertProduct(product)
	if err != nil {
		t.Fatalf("InsertProduct() failed: %s", err.Error())
	}
	defer func() {
		_, _ = catalog.DeleteProductById(created.Product_id)
	}()

	if created.Product_id == 0 {
		t.Fatal("Expected non-zero Product_id after insertion")
	}

	if created.Name != product.Name {
		t.Errorf("Expected Name %q, got %q", product.Name, created.Name)
	}

	retrieved, err := catalog.GetProductById(created.Product_id)
	if err != nil {
		t.Fatalf("GetProductById() failed: %s", err.Error())
	}

	if retrieved.Product_id != created.Product_id {
		t.Errorf("Expected Product_id %d, got %d", created.Product_id, retrieved.Product_id)
	}

	if retrieved.Supplier_id != supplier.Supplier_id {
		t.Errorf("Expected Supplier_id %d, got %d", supplier.Supplier_id, retrieved.Supplier_id)
	}
}

func TestGetProductBySupplierIdAndGetAllProducts(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() failed: %s", err.Error())
	}

	cat, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("newTestCategory() failed: %s", err.Error())
	}

	activeProduct := newTestProduct(supplier.Supplier_id, cat.Category_id)
	inactiveProduct := newTestProduct(supplier.Supplier_id, cat.Category_id)
	inactiveProduct.Is_active = false

	activeCreated, err := catalog.InsertProduct(activeProduct)
	if err != nil {
		t.Fatalf("InsertProduct(active) failed: %s", err.Error())
	}
	defer func() {
		_, _ = catalog.DeleteProductById(activeCreated.Product_id)
	}()

	inactiveCreated, err := catalog.InsertProduct(inactiveProduct)
	if err != nil {
		t.Fatalf("InsertProduct(inactive) failed: %s", err.Error())
	}
	defer func() {
		_, _ = catalog.DeleteProductById(inactiveCreated.Product_id)
	}()

	productsBySupplier, err := catalog.GetProductBySupplierId(supplier.Supplier_id)
	if err != nil {
		t.Fatalf("GetProductBySupplierId() failed: %s", err.Error())
	}

	if len(productsBySupplier) < 2 {
		t.Errorf("Expected at least 2 products for supplier %d, got %d", supplier.Supplier_id, len(productsBySupplier))
	}

	allProducts, err := catalog.GetAllProducts()
	if err != nil {
		t.Fatalf("GetAllProducts() failed: %s", err.Error())
	}

	foundActive := false
	foundInactive := false
	for _, p := range allProducts {
		if p.Product_id == activeCreated.Product_id {
			foundActive = true
		}
		if p.Product_id == inactiveCreated.Product_id {
			foundInactive = true
		}
	}

	if !foundActive {
		t.Errorf("Expected active product %d to appear in GetAllProducts()", activeCreated.Product_id)
	}

	if foundInactive {
		t.Errorf("Did not expect inactive product %d to appear in GetAllProducts()", inactiveCreated.Product_id)
	}
}

func TestUpdateProduct(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() failed: %s", err.Error())
	}

	cat, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("newTestCategory() failed: %s", err.Error())
	}

	product := newTestProduct(supplier.Supplier_id, cat.Category_id)
	created, err := catalog.InsertProduct(product)
	if err != nil {
		t.Fatalf("InsertProduct() failed: %s", err.Error())
	}
	defer func() {
		_, _ = catalog.DeleteProductById(created.Product_id)
	}()

	created.Name = "Updated Product Name"
	created.Description = "Updated description"
	created.Price = 2000
	created.Discount_percent = 10
	created.Stock_quantity = 50
	created.Is_active = true

	rows, err := catalog.UpdateProduct(created)
	if err != nil {
		t.Fatalf("UpdateProduct() failed: %s", err.Error())
	}

	if rows != 1 {
		t.Errorf("Expected 1 row affected, got %d", rows)
	}

	updated, err := catalog.GetProductById(created.Product_id)
	if err != nil {
		t.Fatalf("GetProductById() after update failed: %s", err.Error())
	}

	if updated.Name != created.Name {
		t.Errorf("Expected Name %q, got %q", created.Name, updated.Name)
	}

	if updated.Price != created.Price {
		t.Errorf("Expected Price %d, got %d", created.Price, updated.Price)
	}
}

func TestDeleteProductById(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() failed: %s", err.Error())
	}

	cat, err := db.InsertCategory(newTestCategory())
	if err != nil {
		t.Fatalf("newTestCategory() failed: %s", err.Error())
	}
	fmt.Printf("Created Category with ID: %d\n", cat.Category_id)

	product := newTestProduct(supplier.Supplier_id, cat.Category_id)
	created, err := catalog.InsertProduct(product)
	if err != nil {
		t.Fatalf("InsertProduct() failed: %s", err.Error())
	}

	rows, err := catalog.DeleteProductById(created.Product_id)
	if err != nil {
		t.Fatalf("DeleteProductById() failed: %s", err.Error())
	}

	if rows != 1 {
		t.Errorf("Expected 1 row deleted, got %d", rows)
	}

	_, err = catalog.GetProductById(created.Product_id)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("Expected pgx.ErrNoRows after delete, got %v", err)
	}
}

func TestGetProductByIdNotFound(t *testing.T) {
	setupTestDBConfigAndConnection(t)

	_, err := catalog.GetProductById(-1)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("Expected pgx.ErrNoRows for non-existent product, got %v", err)
	}
}
