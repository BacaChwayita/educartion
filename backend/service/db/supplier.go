package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const supplierColumns = `
	supplier_id,
	name,
	contact_email,
	contact_phone
`

func scanSupplier(row pgx.Row) (model.Supplier, error) {
	var s model.Supplier

	err := row.Scan(
		&s.Supplier_id,
		&s.Name,
		&s.Contact_email,
		&s.Contact_phone,
	)

	return s, err

}

func InsertSupplier(s model.Supplier) (model.Supplier, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO supplier (
		name,
		contact_email,
		contact_phone
	)
	VALUES ($1, $2, $3)
	RETURNING supplier_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		s.Name,
		s.Contact_email,
		s.Contact_phone,
	).Scan(&s.Supplier_id)

	if err == pgx.ErrNoRows {
		return model.Supplier{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertSupplier"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Supplier{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Supplier with ID: %d\n", s.Supplier_id))
	return s, nil
}

func GetSupplierById(supplier_id int) (model.Supplier, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + supplierColumns +
		`
			FROM supplier
			WHERE supplier_id = $1
		`

	var s model.Supplier

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		supplier_id,
	)

	s, err := scanSupplier(row)

	if err == pgx.ErrNoRows {
		return model.Supplier{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetSupplierById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("supplier_id", supplier_id),
		)
		return model.Supplier{}, fmt.Errorf("Error on db select: %w", err)
	}

	return s, nil
}

func UpdateSupplier(s model.Supplier) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE supplier
	Set name = $2,
		contact_email = $3,
		contact_phone = $4

	WHERE supplier_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		s.Supplier_id,
		s.Name,
		s.Contact_email,
		s.Contact_phone,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateSupplier"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("supplier_id", s.Supplier_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteSupplierById(supplier_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM supplier
	WHERE supplier_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		supplier_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteSupplierById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("supplier_id", supplier_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}
