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
	order_id,
	product_id,
	quantity,
	unit_price,
	discount_amount
`

func scanOrderItem(row pgx.Row) (model.Order_item, error) {
	var oi model.Order_item

	err := row.Scan(
		&oi.Order_id,
		&oi.Product_id,
		&oi.Quantity,
		&oi.Unit_price,
		&oi.Discount_amount,
	)

	return oi, err
}

func InsertOrderItem(oi model.Order_item) (model.Order_item, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO order_item (
		order_id,
		product_id,
		quantity,
		unit_price,
		discount_amount
	)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		oi.Order_id,
		oi.Product_id,
		oi.Quantity,
		oi.Unit_price,
		oi.Discount_amount,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertOrderItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Order_item{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created OrderItem with IDs: Order_id(%d), Product_id(%d)\n", oi.Order_id, oi.Product_id))
	return oi, nil
}

func GetOrderItemById(order_id, product_id int) (model.Order_item, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + orderItemColumns +
		`
		FROM order_item
		WHERE order_id   = $1
		  AND product_id = $2
		`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
		product_id,
	)

	oi, err := scanOrderItem(row)

	if err == pgx.ErrNoRows {
		return model.Order_item{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetOrderItemById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
			slog.Int("product_id", product_id),
		)
		return model.Order_item{}, fmt.Errorf("Error on db select: %w", err)
	}

	return oi, nil
}

func UpdateOrderItem(oi model.Order_item) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE order_item
	SET quantity        = $3,
		unit_price      = $4,
		discount_amount = $5
	WHERE order_id   = $1
	  AND product_id = $2
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		oi.Order_id,
		oi.Product_id,
		oi.Quantity,
		oi.Unit_price,
		oi.Discount_amount,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateOrderItem"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", oi.Order_id),
			slog.Int("product_id", oi.Product_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteOrderItemById(order_id, product_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM order_item
	WHERE order_id   = $1
	  AND product_id = $2
	  `

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
		product_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db delete",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteOrderItemById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
			slog.Int("product_id", product_id),
		)
		return -1, fmt.Errorf("Error on db delete: %w", err)
	}

	return result.RowsAffected(), nil
}

func GetOrderItems(order_id int) ([]model.Order_item, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + orderItemColumns +
		`
		FROM order_item
		WHERE order_id = $1
		ORDER BY product_id
	`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetOrderItems"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}
	defer rows.Close()

	orderItems := make([]model.Order_item, 0)
	for rows.Next() {
		var oi model.Order_item
		err := rows.Scan(
			&oi.Order_id,
			&oi.Product_id,
			&oi.Quantity,
			&oi.Unit_price,
			&oi.Discount_amount,
		)
		if err != nil {
			cfg.Logs.Logger.Info(
				"Error on db select",
				slog.String("error", err.Error()),
				slog.String("func", "GetOrderItems"),
				slog.String("timestamp", time.Now().GoString()),
				slog.Int("order_id", order_id),
			)
			return nil, fmt.Errorf("Error on db select: %w", err)
		}

		orderItems = append(orderItems, oi)
	}

	if err = rows.Err(); err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetOrderItems"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}

	return orderItems, nil
}
