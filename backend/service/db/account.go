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
		&acc.Role,
		&acc.Login_attempts,
		&acc.Is_active,
		&acc.Created_at,
	)

	return acc, err

}

func InsertAccount(acc model.Account) (model.Account, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO ACCOUNT (
		full_name,
		email,
		password_hash,
		role,
		login_attempts,
		is_active,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING account_id
	`

	err := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		acc.Full_name,
		acc.Email,
		acc.Password_hash,
		acc.Role,
		0,
		true,
		acc.Created_at,
	).Scan(&acc.Account_id)

	if err == pgx.ErrNoRows {
		return model.Account{}, err
	}

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

func GetAccountById(account_id int) (model.Account, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + accountColumns +
		`
			FROM ACCOUNT
			WHERE account_id = $1
		`

	var acc model.Account

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		account_id,
	)

	acc, err := scanAccount(row)

	if err == pgx.ErrNoRows {
		return model.Account{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetAccountById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("account_id", account_id),
		)
		return model.Account{}, fmt.Errorf("Error on db select: %w", err)
	}

	return acc, nil
}

func UpdateAccount(acc model.Account) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE ACCOUNT
	Set full_name = $2,
		email = $3,
		password_hash = $4,
		role = $5,
		login_attempts = $6,
		is_active = $7,
		created_at = $8

	WHERE account_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		acc.Account_id,
		acc.Full_name,
		acc.Email,
		acc.Password_hash,
		acc.Role,
		acc.Login_attempts,
		acc.Is_active,
		acc.Created_at,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateAccount"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("account_id", acc.Account_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteAccountById(account_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM ACCOUNT
	WHERE account_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		account_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteAccountById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("account_id", account_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func GetAccountByJTI(token_string string) (model.Account, error) {
	cfg := config.GetConfig()

	sql := `
	SELECT 
		a.account_id,
		a.full_name,
		a.email,
		a.password_hash,
		a.role,
		a.login_attempts,
		a.is_active,
		a.created_at
	FROM ACCOUNT a
	INNER JOIN AccountLogin al
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
		&acc.Role,
		&acc.Login_attempts,
		&acc.Is_active,
		&acc.Created_at,
	)

	if err == pgx.ErrNoRows {
		return model.Account{}, err
	}

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

func GetAccountByEmail(email string) (model.Account, error) {
	cfg := config.GetConfig()

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

	if err == pgx.ErrNoRows {
		return model.Account{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetUserByEmail"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("email", email),
		)
		return model.Account{}, fmt.Errorf("Error on db select: %w", err)
	}

	return acc, nil
}

func IncLoginAttempts(account_id int) (int64, error) {
	cfg := config.GetConfig()

	// Doing funny case statement since there's no easy select account values yet,
	// doesn't make sense to add to parms (db should be source of truth),
	// and I don't want to do 2 sql statements if we can do it in 1
	sql := `
	UPDATE ACCOUNT
	SET login_attempts = login_attempts + 1,
	is_active = case when login_attempts + 1 >= 3 then false else is_active end
	WHERE account_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
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
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil

}

func InsertAccountLogin(al model.AccountLogin) (model.AccountLogin, error) {
	cfg := config.GetConfig()

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

	if err == pgx.ErrNoRows {
		return model.AccountLogin{}, err
	}

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

func DeleteAccountLoginById(token_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM ACCOUNTLOGIN
	WHERE token_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		token_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteAccountLoginById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("token_id", token_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func DeleteAccountLoginByToken(token string) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM ACCOUNTLOGIN
	WHERE token_string = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		token,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteAccountLoginById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("token", token),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}
