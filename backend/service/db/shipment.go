package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const shipmentColumns = `
		shipment_id,
		order_id,
		courier_name,
		courier_type,
		delivery_reference,
		delivery_address_line1,
		delivery_address_line2,
		delivery_city,
		delivery_state,
		delivery_postal_code,
		delivery_country,
		recipient_name,
		recipient_phone,
		shipment_status,
		dispatched_at,
		delivered_at,
		created_at
`

func scanShipment(row pgx.Row) (model.Shipment, error) {
	var s model.Shipment

	err := row.Scan(
		&s.Shipment_id,
		&s.Order_id,
		&s.Courier_name,
		&s.Courier_type,
		&s.Delivery_reference,
		&s.Delivery_address_line1,
		&s.Delivery_address_line2,
		&s.Delivery_city,
		&s.Delivery_state,
		&s.Delivery_postal_code,
		&s.Delivery_country,
		&s.Recipient_name,
		&s.Recipient_phone,
		&s.Shipment_status,
		&s.Dispatched_at,
		&s.Delivered_at,
		&s.Created_at,
	)

	return s, err

}

func InsertShipment(s model.Shipment) (model.Shipment, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO shipment (
		order_id,
		courier_name,
		courier_type,
		delivery_reference,
		delivery_address_line1,
		delivery_address_line2,
		delivery_city,
		delivery_state,
		delivery_postal_code,
		delivery_country,
		recipient_name,
		recipient_phone,
		shipment_status,
		dispatched_at,
		delivered_at,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,$13, $14, $15, $16)
	RETURNING shipment_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		s.Order_id,
		s.Courier_name,
		s.Courier_type,
		s.Delivery_reference,
		s.Delivery_address_line1,
		s.Delivery_address_line2,
		s.Delivery_city,
		s.Delivery_state,
		s.Delivery_postal_code,
		s.Delivery_country,
		s.Recipient_name,
		s.Recipient_phone,
		s.Shipment_status,
		s.Dispatched_at,
		s.Delivered_at,
		s.Created_at,
	).Scan(&s.Shipment_id)

	if err == pgx.ErrNoRows {
		return model.Shipment{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertShipment"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", s.Order_id),
		)
		return model.Shipment{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Shipment with ID: %d\n", s.Shipment_id))
	return s, nil
}

func GetShipmentById(shipment_id int) (model.Shipment, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + shipmentColumns +
		`
			FROM shipment
			WHERE shipment_id = $1
		`

	var s model.Shipment

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		shipment_id,
	)

	s, err := scanShipment(row)

	if err == pgx.ErrNoRows {
		return model.Shipment{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetShipmentById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("shipment_id", shipment_id),
		)
		return model.Shipment{}, fmt.Errorf("Error on db select: %w", err)
	}

	return s, nil
}

func GetShipmentByOrderId(order_id int) (model.Shipment, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + shipmentColumns +
		`
			FROM shipment
			WHERE order_id = $1
		`

	var s model.Shipment

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		order_id,
	)

	s, err := scanShipment(row)

	if err == pgx.ErrNoRows {
		return model.Shipment{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetShipmentByOrderId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("order_id", order_id),
		)
		return model.Shipment{}, fmt.Errorf("Error on db select: %w", err)
	}

	return s, nil
}
func UpdateShipment(s model.Shipment) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE shipment
	SET order_id = $2,
		courier_name = $3,
		courier_type = $4,
		delivery_reference = $5,
		delivery_address_line1 = $6,
		delivery_address_line2 = $7,
		delivery_city = $8,
		delivery_state = $9,
		delivery_postal_code = $10,
		delivery_country = $11,
		recipient_name = $12,
		recipient_phone = $13,
		shipment_status = $14,
		dispatched_at = $15,
		delivered_at = $16,
		created_at = $17

	WHERE shipment_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		s.Shipment_id,
		s.Order_id,
		s.Courier_name,
		s.Courier_type,
		s.Delivery_reference,
		s.Delivery_address_line1,
		s.Delivery_address_line2,
		s.Delivery_city,
		s.Delivery_state,
		s.Delivery_postal_code,
		s.Delivery_country,
		s.Recipient_name,
		s.Recipient_phone,
		s.Shipment_status,
		s.Dispatched_at,
		s.Delivered_at,
		s.Created_at,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateShipment"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("shipment_id", s.Shipment_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteShipmentById(shipment_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM shipment
	WHERE shipment_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		shipment_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteShipmentById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("shipment_id", shipment_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}
