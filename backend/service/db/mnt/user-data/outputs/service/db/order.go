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
	supplier_id,
	status,
	notes,
	created_at,
	updated_at
`

func scanOrder(row pgx.Row) (model.Order, error) {
	var o model.Order

	err := row.Scan(
		&o.Order_id,
		&o.Account_id,
		&o.Supplier_id,
		&o.Status,
		&o.Notes,
		&o.Created_at,
		&o.Updated_at,
	)

	return o, err
}

func InsertOrder(o model.Order) (model.Order, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO orders (
		account_id,
		supplier_id,
		status,
		notes,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING order_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		o.Account_id,
		o.Supplier_id,
		o.Status,
		o.Notes,
		o.Created_at,
		o.Updated_at,
	).Scan(&o.Order_id)

	if err == pgx.ErrNoRows {
		return model.Order{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertOrder"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Order{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Order with ID: %d\n", o.Order_id))
	return o, nil
}

func GetOrderById(order_id int) (model.Order, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + orderColumns +
		`
		FROM orders
		WHERE order_id = $1
		`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
	)

	o, err := scanOrder(row)

	if err == pgx.ErrNoRows {
		return model.Order{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetOrderById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return model.Order{}, fmt.Errorf("Error on db select: %w", err)
	}

	return o, nil
}

func UpdateOrder(o model.Order) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE orders
	SET account_id  = $2,
		supplier_id = $3,
		status      = $4,
		notes       = $5,
		updated_at  = $6
	WHERE order_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		o.Order_id,
		o.Account_id,
		o.Supplier_id,
		o.Status,
		o.Notes,
		o.Updated_at,
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
			"Error on db delete",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteOrderById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return -1, fmt.Errorf("Error on db delete: %w", err)
	}

	return result.RowsAffected(), nil
}
