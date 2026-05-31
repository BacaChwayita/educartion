package db_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func insertTestProductImage(t *testing.T, productID int, imageURL string, isPrimary bool, sortOrder int) int {
	t.Helper()

	cfg := config.GetConfig()
	var imageID int
	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		`INSERT INTO product_image (product_id, image_url, alt_text, is_primary, sort_order)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING product_image_id`,
		productID,
		imageURL,
		"test image",
		isPrimary,
		sortOrder,
	).Scan(&imageID)
	if err != nil {
		t.Fatalf("Failed to insert product image: %s", err.Error())
	}

	return imageID
}

func TestGetProductImagesByProductId(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	product, err := db.InsertProduct(newTestProduct(supplier.Supplier_id))
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	primaryURL := "primary-test.jpg"
	secondaryURL := "secondary-test.jpg"
	insertTestProductImage(t, product.Product_id, secondaryURL, false, 2)
	insertTestProductImage(t, product.Product_id, primaryURL, true, 1)

	images, err := db.GetProductImagesByProductId(product.Product_id)
	if err != nil {
		t.Fatalf("GetProductImagesByProductId() call failed: %s", err.Error())
	}

	if len(images) != 2 {
		t.Fatalf("Expected 2 images, got %d", len(images))
	}
	if !images[0].Is_primary {
		t.Fatalf("Expected first image to be primary")
	}
	if images[0].Image_url != primaryURL {
		t.Fatalf("Expected primary image URL %q, got %q", primaryURL, images[0].Image_url)
	}
}

func TestGetPrimaryProductImageByProductId(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	product, err := db.InsertProduct(newTestProduct(supplier.Supplier_id))
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	primaryURL := "primary-only.jpg"
	insertTestProductImage(t, product.Product_id, "secondary.jpg", false, 2)
	insertTestProductImage(t, product.Product_id, primaryURL, true, 3)

	image, err := db.GetPrimaryProductImageByProductId(product.Product_id)
	if err != nil {
		t.Fatalf("GetPrimaryProductImageByProductId() call failed: %s", err.Error())
	}
	if !image.Is_primary {
		t.Fatalf("Expected primary image, got non-primary record")
	}
	if image.Image_url != primaryURL {
		t.Fatalf("Expected image URL %q, got %q", primaryURL, image.Image_url)
	}
}

func TestGetPrimaryProductImageByProductIdNoRows(t *testing.T) {
	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	supplier, err := db.InsertSupplier(newTestSupplier())
	if err != nil {
		t.Fatalf("InsertSupplier() call failed: %s", err.Error())
	}

	product, err := db.InsertProduct(newTestProduct(supplier.Supplier_id))
	if err != nil {
		t.Fatalf("InsertProduct() call failed: %s", err.Error())
	}

	_, err = db.GetPrimaryProductImageByProductId(product.Product_id)
	if err != pgx.ErrNoRows {
		t.Fatalf("Expected pgx.ErrNoRows, got %v", err)
	}
}
