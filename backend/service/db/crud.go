package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
)

func InsertAccount(cfg *config.Config, acc model.Account) (model.Account, error) {
	sql := `
	INSERT INTO ACCOUNT (
		full_name,
		email,
		password_hash,
		password_salt,
		role,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING account_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		acc.Full_name,
		acc.Email,
		acc.Password_hash,
		acc.Password_salt,
		acc.Role,
		acc.Created_at,
	).Scan(&acc.Account_id)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertAccount"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Account{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Account with ID: %d\n", acc.Account_id))
	return acc, nil
}
