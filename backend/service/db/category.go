package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const categoryColumns = `
	category_id,
	name,
	description,
	is_active
`

func scanCategory(row pgx.Row) (model.Category, error) {
	var c model.Category

	err := row.Scan(
		&c.Category_id,
		&c.Name,
		&c.Description,
		&c.Is_active,
	)

	return c, err
}

func InsertCategory(c model.Category) (model.Category, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO category (
		name,
		description,
		is_active
	)
	VALUES ($1, $2, $3)
	RETURNING category_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		c.Name,
		c.Description,
		c.Is_active,
	).Scan(&c.Category_id)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertCategory"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Category{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Category with ID: %d", c.Category_id))
	return c, nil
}

func GetCategoryById(category_id int) (model.Category, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + categoryColumns + `
	FROM category
	WHERE category_id = $1
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		category_id,
	)

	return scanCategory(row)
}

func GetAllCategories() ([]model.Category, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + categoryColumns + `
	FROM category
	ORDER BY category_id
	`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetAllCategories"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		err := rows.Scan(
			&c.Category_id,
			&c.Name,
			&c.Description,
			&c.Is_active,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		cfg.Logs.Logger.Info(
			"Error iterating category rows",
			slog.String("error", err.Error()),
			slog.String("func", "GetAllCategories"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return nil, fmt.Errorf("Error iterating category rows: %w", err)
	}

	return categories, nil
}

func UpdateCategory(c model.Category) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE category
	SET name = $2,
		description = $3,
		is_active = $4
	WHERE category_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		c.Category_id,
		c.Name,
		c.Description,
		c.Is_active,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateCategory"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("category_id", c.Category_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteCategoryById(category_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM category
	WHERE category_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		category_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db delete",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteCategoryById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("category_id", category_id),
		)
		return -1, fmt.Errorf("Error on db delete: %w", err)
	}

	return result.RowsAffected(), nil
}

func GetProductsByCategoryId(category_id int) ([]model.Product, error) {
	cfg := config.GetConfig()

	sql := `
	SELECT product_id, supplier_id, category_id, name, description, price, discount_percent, stock_quantity, is_active
	FROM product
	WHERE category_id = $1
	AND is_active = true
	`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
		category_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetProductsByCategoryId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("category_id", category_id),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(
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
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		cfg.Logs.Logger.Info(
			"Error iterating product rows",
			slog.String("error", err.Error()),
			slog.String("func", "GetProductsByCategoryId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("category_id", category_id),
		)
		return nil, fmt.Errorf("Error iterating product rows: %w", err)
	}

	return products, nil
}
