package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const accountColumns = `
	account_id,
	full_name,
	email,
	password_hash,
	password_salt,
	role,
	login_attempts,
	is_active,
	created_at
`

func scanAccount(row pgx.Row) (model.Account, error) {
	var acc model.Account

	err := row.Scan(
		&acc.Account_id,
		&acc.Full_name,
		&acc.Email,
		&acc.Password_hash,
		&acc.Password_salt,
		&acc.Role,
		&acc.Login_attempts,
		&acc.Is_active,
		&acc.Created_at,
	)

	return acc, err

}

func InsertAccount(cfg *config.Config, acc model.Account) (model.Account, error) {
	sql := `
	INSERT INTO ACCOUNT (
		full_name,
		email,
		password_hash,
		password_salt,
		role,
		login_attempts,
		is_active,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
		0,
		true,
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

func GetAccountWithJTI(cfg *config.Config, token_string string) (model.Account, error) {
	sql := `
	SELECT 
		a.account_id,
		a.full_name,
		a.email,
		a.password_hash,
		a.password_salt,
		a.role,
		a.login_attempts,
		a.is_active,
		a.created_at,
	FROM ACCOUNT a
	INNER JOIN AccountLogins al
	  ON al.account_id = a.account_id
	WHERE al.token_string = $1
	`

	var acc model.Account

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		token_string,
	).Scan(
		&acc.Account_id,
		&acc.Full_name,
		&acc.Email,
		&acc.Password_hash,
		&acc.Password_salt,
		&acc.Role,
		&acc.Login_attempts,
		&acc.Is_active,
		&acc.Created_at,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetAccountWithJTI"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Account{}, fmt.Errorf("Error on db select: %w", err)
	}

	return acc, nil
}

func GetUserByEmail(cfg *config.Config, email string) (model.Account, error) {
	// TODO: Probably need a getAccount() that will do all getting, and we just pass the values that needs to be filtered on, to ensure we don't copy paste this same select with a slightly different where clause (imagine adding a column to the table, have to add all over...)
	sql := `SELECT ` + accountColumns +
		`
			FROM ACCOUNT
			WHERE email = $1
		`

	var acc model.Account

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		email,
	)

	acc, err := scanAccount(row)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "getUserByEmail"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("email", email),
		)
		return model.Account{}, fmt.Errorf("Error on db select: %w", err)
	}

	return acc, nil
}

func IncLoginAttempts(cfg *config.Config, account_id int) error {
	// Doing funny case statement since there's no easy select account values yet,
	// doesn't make sense to add to parms (db should be source of truth),
	// and I don't want to do 2 sql statements if we can do it in 1
	sql := `
	UPDATE ACCOUNT
	SET login_attempts = login_attempts + 1,
	is_active = case when login_attempts + 1 >= 3 then false else is_active end
	WHERE account_id = $1
	`

	// TODO: Look into the command tag returned...
	_, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		account_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "incLoginAttempts"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("account_id", account_id),
		)
		return fmt.Errorf("Error on db update: %w", err)
	}

	return nil

}

func InsertAccountLogin(cfg *config.Config, al model.AccountLogin) (model.AccountLogin, error) {
	sql := `
	INSERT INTO ACCOUNTLOGIN (
		account_id,
		token_string
	)
	VALUES ($1, $2)
	RETURNING token_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		al.Account_id,
		al.Token_string,
	).Scan(&al.Token_id)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertAccountLogin"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.AccountLogin{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Account Login with ID: %d\n", al.Account_id))
	return al, nil
}
