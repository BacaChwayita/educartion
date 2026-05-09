package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const orderColumns = `
		order_id,
		account_id,
		order_number,
		status,
		subtotal_amount,
		discount_amount,
		total_amount,
		placed_at
`

func scanOrder(row pgx.Row) (model.Orders, error) {
	var o model.Orders

	err := row.Scan(
		&o.Order_id,
		&o.Account_id,
		&o.Order_number,
		&o.Status,
		&o.Subtotal_amount,
		&o.Discount_amount,
		&o.Total_amount,
		&o.Placed_at,
	)

	return o, err

}

func InsertOrder(o model.Orders) (model.Orders, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO orders (
		account_id,
		order_number,
		status,
		subtotal_amount,
		discount_amount,
		total_amount,
		placed_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING order_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		o.Account_id,
		o.Order_number,
		o.Status,
		o.Subtotal_amount,
		o.Discount_amount,
		o.Total_amount,
		o.Placed_at,
	).Scan(&o.Order_id)

	if err == pgx.ErrNoRows {
		return model.Orders{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertOrder"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Orders{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Order with ID: %d\n", o.Order_id))
	return o, nil
}

func GetOrderById(order_id int) (model.Orders, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + orderColumns +
		`
			FROM orders
			WHERE order_id = $1
		`

	var o model.Orders

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
	)

	o, err := scanOrder(row)

	if err == pgx.ErrNoRows {
		return model.Orders{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetOrderById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return model.Orders{}, fmt.Errorf("Error on db select: %w", err)
	}

	return o, nil
}

func UpdateOrder(o model.Orders) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE orders
	Set account_id = $2,
		order_number = $3,
		status = $4,
		subtotal_amount = $5,
		discount_amount = $6,
		total_amount = $7,
		placed_at = $8

	WHERE order_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		o.Order_id,
		o.Account_id,
		o.Order_number,
		o.Status,
		o.Subtotal_amount,
		o.Discount_amount,
		o.Total_amount,
		o.Placed_at,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateOrder"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", o.Order_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteOrderById(order_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM orders
	WHERE order_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteOrderById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}
