package db

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const productColumns = `
	product_id,
	supplier_id,
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
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Discount_percent,
		&p.Stock_quantity,
		&p.Is_active,
	)

	return p, err
}

func InsertProduct(p model.Product) (model.Product, error) {

	cfg := config.GetConfig()

	sql := `
	INSERT INTO product (
		supplier_id,
		name,
		description,
		price,
		discount_percent,
		stock_quantity,
		is_active
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING product_id

	`
	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		p.Supplier_id,
		p.Name,
		p.Description,
		p.Price,
		p.Discount_percent,
		p.Stock_quantity,
		p.Is_active,
	).Scan(&p.Product_id)

	if err == pgx.ErrNoRows {
		if err != nil {

			cfg.Logs.Logger.Info(
				"Error on db insert",
				slog.String("error", err.Error()),
				slog.String("func", "InsertProduct"),
				slog.String("timestamp", time.Now().GoString()),
			)

			return model.Product{}, fmt.Errorf("Error on db insert: %w", err)
		}

	}

	cfg.Logs.Logger.Info(
		fmt.Sprintf("Created Product with ID: %d\n", p.Product_id),
	)

	return p, nil
}

func GetProductById(product_id int) (model.Product, error) {

	cfg := config.GetConfig()

	sql := `SELECT ` + productColumns +
		`
		FROM product
		WHERE product_id = $1
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		product_id,
	)

	p, err := scanProduct(row)

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

	return p, nil
}

func GetAllProducts() ([]model.Product, error) {

	cfg := config.GetConfig()

	sql := `
	SELECT ` + productColumns + `
	FROM product
	ORDER BY created_at DESC
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

		return nil, err
	}

	defer rows.Close()

	var products []model.Product
	for rows.Next() {

		var p model.Product

		err := rows.Scan(
			&p.Product_id,
			&p.Supplier_id,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Discount_percent,
			&p.Stock_quantity,
			&p.Is_active,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func UpdateProduct(p model.Product) (int64, error) {

	cfg := config.GetConfig()

	sql := `
	UPDATE product
	SET supplier_id = $2,
		name = $3,
		description = $4,
		price = $5,
		discount_percent = $6,
		stock_quantity = $7,
		is_active = $8

	WHERE product_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		p.Product_id,
		p.Supplier_id,
		p.Name,
		p.Description,
		p.Price,
		p.Discount_percent,
		p.Stock_quantity,
		p.Is_active,
	)

	if err != nil {

		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateProduct"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", p.Product_id),
		)

		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteProductById(product_id int) (int64, error) {

	cfg := config.GetConfig()

	sql := `
	DELETE FROM product
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

func SearchProduct(products []model.Product, query string) ([]model.Product, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	var results []model.Product

	for _, p := range products {
		// simple case-insensitive contains search
		if containsIgnoreCase(p.Name, query) {
			results = append(results, p)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no products found")
	}

	return results, nil
}

func containsIgnoreCase(input string, query string) bool {
	return strings.Contains(strings.ToLower(input), strings.ToLower(query))
}

type ProductSearchCriteria struct {
	Search_name string
	Supplier_id int
	Min_price   int
	Max_price   int
}

func GetProductsWithSearch(psc ProductSearchCriteria) ([]model.Product, error) {
	products, err := GetAllProducts()
	if err != nil {
		return nil, err
	}

	if psc.Search_name == "" && psc.Supplier_id == 0 && psc.Min_price == 0 && psc.Max_price == 0 {
		return products, nil
	}

	results := make([]model.Product, 0, len(products))
	for _, p := range products {
		if psc.Search_name != "" && !containsIgnoreCase(p.Name, psc.Search_name) {
			continue
		}
		if psc.Supplier_id != 0 && p.Supplier_id != psc.Supplier_id {
			continue
		}
		if psc.Min_price != 0 && p.Price < psc.Min_price {
			continue
		}
		if psc.Max_price != 0 && p.Price > psc.Max_price {
			continue
		}
		results = append(results, p)
	}

	if len(results) == 0 {
		return nil, errors.New("no products found")
	}

	return results, nil
}
