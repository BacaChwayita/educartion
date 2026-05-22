package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const paymentColumns = `
		payment_id,
		order_id,
		payment_method,
		payment_status,
		transaction_ref,
		amount,
		paid_at
`

func scanPayment(row pgx.Row) (model.Payment, error) {
	var p model.Payment

	err := row.Scan(
		&p.Payment_id,
		&p.Order_id,
		&p.Payment_method,
		&p.Payment_status,
		&p.Transaction_ref,
		&p.Amount,
		&p.Paid_at,
	)

	return p, err
}

func InsertPayment(p model.Payment) (model.Payment, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO payment (
		order_id,
		payment_method,
		payment_status,
		transaction_ref,
		amount,
		paid_at
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING payment_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		p.Order_id,
		p.Payment_method,
		p.Payment_status,
		p.Transaction_ref,
		p.Amount,
		p.Paid_at,
	).Scan(&p.Payment_id)

	if err == pgx.ErrNoRows {
		return model.Payment{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertPayment"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", p.Order_id),
		)
		return model.Payment{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Payment with ID: %d\n", p.Payment_id))
	return p, nil
}

func GetPaymentById(payment_id int) (model.Payment, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + paymentColumns + `
		FROM payment
		WHERE payment_id = $1
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		payment_id,
	)

	p, err := scanPayment(row)

	if err == pgx.ErrNoRows {
		return model.Payment{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetPaymentById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("payment_id", payment_id),
		)
		return model.Payment{}, fmt.Errorf("Error on db select: %w", err)
	}

	return p, nil
}

func GetPaymentByOrderId(order_id int) (model.Payment, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + paymentColumns + `
		FROM payment
		WHERE order_id = $1
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
	)

	p, err := scanPayment(row)

	if err == pgx.ErrNoRows {
		return model.Payment{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetPaymentByOrderId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return model.Payment{}, fmt.Errorf("Error on db select: %w", err)
	}

	return p, nil
}
