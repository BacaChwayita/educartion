package catalog

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const productColumns = `
	product_id,
	supplier_id,
	category_id,
	name,
	description,
	price,
	discount_percent,
	stock_quantity,
	is_active
`

func scanProduct(row pgx.Row) (model.Product, error) {
	var p model.Product

	err := row.Scan(
		&p.Product_id,
		&p.Supplier_id,
		&p.Category_id,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Discount_percent,
		&p.Stock_quantity,
		&p.Is_active,
	)

	return p, err
}

func InsertProduct(product model.Product) (model.Product, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO PRODUCT (
		supplier_id,
		category_id,
		name,
		description,
		price,
		discount_percent,
		stock_quantity,
		is_active
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING ` + productColumns

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		product.Supplier_id,
		product.Category_id,
		product.Name,
		product.Description,
		product.Price,
		product.Discount_percent,
		product.Stock_quantity,
		product.Is_active,
	)

	newProduct, err := scanProduct(row)

	if err == pgx.ErrNoRows {
		return model.Product{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertProduct"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Product{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Product with ID: %d\n", newProduct.Product_id))
	return newProduct, nil
}

func GetProductById(product_id int) (model.Product, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + productColumns +
		`
			FROM PRODUCT
			WHERE product_id = $1
		`

	var product model.Product

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		product_id,
	)

	product, err := scanProduct(row)

	if err == pgx.ErrNoRows {
		return model.Product{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetProductById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", product_id),
		)
		return model.Product{}, fmt.Errorf("Error on db select: %w", err)
	}

	return product, nil
}

func GetProductBySupplierId(supplier_id int) ([]model.Product, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + productColumns +
		`
			FROM PRODUCT
			WHERE supplier_id = $1
		`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
		supplier_id,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetProductBySupplierId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("supplier_id", supplier_id),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			cfg.Logs.Logger.Info(
				"Error scanning product",
				slog.String("error", err.Error()),
				slog.String("func", "GetProductBySupplierId"),
				slog.String("timestamp", time.Now().GoString()),
			)
			return nil, fmt.Errorf("Error scanning product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func GetAllProducts() ([]model.Product, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + productColumns +
		`
			FROM PRODUCT
			WHERE is_active = true
		`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetAllProducts"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			cfg.Logs.Logger.Info(
				"Error scanning product",
				slog.String("error", err.Error()),
				slog.String("func", "GetAllProducts"),
				slog.String("timestamp", time.Now().GoString()),
			)
			return nil, fmt.Errorf("Error scanning product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func UpdateProduct(product model.Product) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE PRODUCT
	SET
		supplier_id = $2,
		category_id = $3,
		name = $4,
		description = $5,
		price = $6,
		discount_percent = $7,
		stock_quantity = $8,
		is_active = $9
	WHERE product_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		product.Product_id,
		product.Supplier_id,
		product.Category_id,
		product.Name,
		product.Description,
		product.Price,
		product.Discount_percent,
		product.Stock_quantity,
		product.Is_active,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateProduct"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", product.Product_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	if result.RowsAffected() == 0 {
		cfg.Logs.Logger.Info(
			"Product not found for update",
			slog.Int("product_id", product.Product_id),
			slog.String("func", "UpdateProduct"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return 0, pgx.ErrNoRows
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Updated Product with ID: %d\n", product.Product_id))
	return result.RowsAffected(), nil
}

func DeleteProductById(product_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM PRODUCT
	WHERE product_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		product_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db delete",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteProductById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", product_id),
		)
		return -1, fmt.Errorf("Error on db delete: %w", err)
	}

	return result.RowsAffected(), nil
}
