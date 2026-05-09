package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const cartColumns = `
	cart_id,
	account_id,
	session_key,
	total_price,
	created_at,
	updated_at
`

func scanCart(row pgx.Row) (model.Cart, error) {
	var c model.Cart

	err := row.Scan(
		&c.Cart_id,
		&c.Account_id,
		&c.Session_key,
		&c.Total_Price,
		&c.Created_at,
		&c.Updated_at,
	)

	return c, err
}

func InsertCart(cart model.Cart) (model.Cart, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO CART (
		account_id,
		session_key,
		total_price,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING ` + cartColumns

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		cart.Account_id,
		cart.Session_key,
		cart.Total_Price,
		cart.Created_at,
		cart.Updated_at,
	)

	newCart, err := scanCart(row)

	if err == pgx.ErrNoRows {
		return model.Cart{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertCart"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart{}, fmt.Errorf("Error on db insert: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Created Cart with ID: %d\n", newCart.Cart_id))
	return newCart, nil
}

func GetCartById(cart_id int) (model.Cart, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + cartColumns +
		`
			FROM CART
			WHERE cart_id = $1
		`

	var cart model.Cart

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		cart_id,
	)

	cart, err := scanCart(row)

	if err == pgx.ErrNoRows {
		return model.Cart{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("cart_id", cart_id),
		)
		return model.Cart{}, fmt.Errorf("Error on db select: %w", err)
	}

	return cart, nil
}

func GetCartByAccountId(account_id int) (model.Cart, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + cartColumns +
		`
			FROM CART
			WHERE account_id = $1
		`

	var cart model.Cart

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		account_id,
	)

	cart, err := scanCart(row)

	if err == pgx.ErrNoRows {
		return model.Cart{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartByAccountId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("account_id", account_id),
		)
		return model.Cart{}, fmt.Errorf("Error on db select: %w", err)
	}

	return cart, nil
}

func GetCartBySessionKey(session_key string) (model.Cart, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + cartColumns +
		`
			FROM CART
			WHERE session_key = $1
		`

	var cart model.Cart

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		session_key,
	)

	cart, err := scanCart(row)

	if err == pgx.ErrNoRows {
		return model.Cart{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartBySessionKey"),
			slog.String("timestamp", time.Now().GoString()),
			slog.String("session_key", session_key),
		)
		return model.Cart{}, fmt.Errorf("Error on db select: %w", err)
	}

	return cart, nil
}

func UpdateCart(cart model.Cart) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE CART
	SET
		session_key = $2,
		total_price = $3,
		updated_at  = $4
	WHERE cart_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		cart.Cart_id,
		cart.Session_key,
		cart.Total_Price,
		cart.Updated_at,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateCart"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("cart_id", cart.Cart_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	if result.RowsAffected() == 0 {
		cfg.Logs.Logger.Info(
			"Cart not found for update",
			slog.Int("cart_id", cart.Cart_id),
			slog.String("func", "UpdateCart"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return 0, pgx.ErrNoRows
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Updated Cart with ID: %d\n", cart.Cart_id))
	return result.RowsAffected(), nil
}

func DeleteCartById(cart_id int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM CART
	WHERE cart_id = $1
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		cart_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteCartById"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("cart_id", cart_id),
		)
		return -1, fmt.Errorf("Error on db update: %w", err)
	}

	return result.RowsAffected(), nil
}

func GetCartByJTI(token_string string) (model.Cart, error) {
	cfg := config.GetConfig()

	sql := `
	SELECT
		c.cart_id,
		c.account_id,
		c.session_key,
		c.total_price,
		c.created_at,
		c.updated_at
	FROM CART c
	INNER JOIN ACCOUNT a
	  ON a.account_id = c.account_id
	INNER JOIN ACCOUNTLOGIN al
	  ON al.account_id = a.account_id
	WHERE al.token_string = $1
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		token_string,
	)

	cart, err := scanCart(row)

	if err == pgx.ErrNoRows {
		return model.Cart{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartByJTI"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart{}, fmt.Errorf("Error on db select: %w", err)
	}

	return cart, nil
}
