package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const orderItemColumns = `
	order_item_id,
	order_id,
	product_name,
	quantity,
	unit_price
`

func scanOrderItem(row pgx.Row) (model.OrderItem, error) {
	var oi model.OrderItem

	err := row.Scan(
		&oi.Order_item_id,
		&oi.Order_id,
		&oi.Product_name,
		&oi.Quantity,
		&oi.Unit_price,
	)

	return oi, err
}

func InsertOrderItem(oi model.OrderItem) (model.OrderItem, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO order_items (
		order_id,
		product_name,
		quantity,
		unit_price
	)
	VALUES ($1, $2, $3, $4)
	RETURNING order_item_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		oi.Order_id,
		oi.Product_name,
		oi.Quantity,
		oi.Unit_price,
	).Scan(&oi.Order_item_id)

	if err == pgx.ErrNoRows {
		return model.OrderItem{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertOrderItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.OrderItem{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created OrderItem with ID: %d\n", oi.Order_item_id))
	return oi, nil
}

func GetOrderItemById(order_item_id int) (model.OrderItem, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + orderItemColumns +
		`
		FROM order_items
		WHERE order_item_id = $1
		`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		order_item_id,
	)

	oi, err := scanOrderItem(row)

	if err == pgx.ErrNoRows {
		return model.OrderItem{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetOrderItemById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_item_id", order_item_id),
		)
		return model.OrderItem{}, fmt.Errorf("Error on db select: %w", err)
	}

	return oi, nil
}

func UpdateOrderItem(oi model.OrderItem) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE order_items
	SET order_id     = $2,
		product_name = $3,
		quantity     = $4,
		unit_price   = $5
	WHERE order_item_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		oi.Order_item_id,
		oi.Order_id,
		oi.Product_name,
		oi.Quantity,
		oi.Unit_price,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateOrderItem"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_item_id", oi.Order_item_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteOrderItemById(order_item_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM order_items
	WHERE order_item_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		order_item_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db delete",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteOrderItemById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_item_id", order_item_id),
		)
		return -1, fmt.Errorf("Error on db delete: %w", err)
	}

	return result.RowsAffected(), nil
}
