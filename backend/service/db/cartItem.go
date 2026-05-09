package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const cartItemColumns = `
	cart_id,
	product_id,
	quantity,
	unit_price
`

func scanCartItem(row pgx.Row) (model.Cart_item, error) {
	var ci model.Cart_item

	err := row.Scan(
		&ci.Cart_id,
		&ci.Product_id,
		&ci.Quantity,
		&ci.Unit_price,
	)

	return ci, err

}

func InsertCartItem(cartItem model.Cart_item) (model.Cart_item, error) {
	cfg := config.GetConfig()

	sql := `
	INSERT INTO CART_ITEM (
		cart_id,
		product_id,
		quantity,
		unit_price
	)
	VALUES ($1, $2, $3, $4)
	RETURNING ` + cartItemColumns

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		cartItem.Cart_id,
		cartItem.Product_id,
		cartItem.Quantity,
		cartItem.Unit_price,
	)

	newItem, err := scanCartItem(row)

	if err == pgx.ErrNoRows {
		return model.Cart_item{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db insert",
			slog.String("error", err.Error()),
			slog.String("func", "InsertCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart_item{}, fmt.Errorf("Error on db insert: %w", err)
	}
	cfg.Logs.Logger.Info(fmt.Sprintf("Created Cart with ID: %d\n", cartItem.Cart_id))
	return newItem, nil
}

func GetCartItem(cartID int, productID int) (model.Cart_item, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + cartItemColumns +
		`
		FROM CART_ITEM
		WHERE cart_id    = $1
		AND   product_id = $2
		`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		cartID,
		productID,
	)

	cartItem, err := scanCartItem(row)

	if err == pgx.ErrNoRows {
		cfg.Logs.Logger.Info(
			"Cart item not found",
			slog.Int("cart_id", cartID),
			slog.Int("product_id", productID),
			slog.String("func", "GetCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart_item{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart_item{}, fmt.Errorf("Error on db select: %w", err)
	}

	return cartItem, nil
}

func GetCartItemsByCart(cartID int) ([]model.Cart_item, error) {
	cfg := config.GetConfig()

	sql := `
	SELECT
		cart_id,
		product_id,
		quantity,
		unit_price
	FROM CART_ITEM
	WHERE cart_id = $1
	`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
		cartID,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db query",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartItemsByCart"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return nil, fmt.Errorf("Error on db query: %w", err)
	}
	defer rows.Close()

	var items []model.Cart_item

	for rows.Next() {
		var cartItem model.Cart_item
		err := rows.Scan(
			&cartItem.Cart_id,
			&cartItem.Product_id,
			&cartItem.Quantity,
			&cartItem.Unit_price,
		)
		if err != nil {
			cfg.Logs.Logger.Info(
				"Error scanning cart item row",
				slog.String("error", err.Error()),
				slog.String("func", "GetCartItemsByCart"),
				slog.String("timestamp", time.Now().GoString()),
			)
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		items = append(items, cartItem)
	}

	if err = rows.Err(); err != nil {
		cfg.Logs.Logger.Info(
			"Row iteration error",
			slog.String("error", err.Error()),
			slog.String("func", "GetCartItemsByCart"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return nil, fmt.Errorf("Row iteration error: %w", err)
	}

	return items, nil
}

func UpdateCartItem(cartID int, productID int, quantity int) (model.Cart_item, error) {
	cfg := config.GetConfig()

	sql := `
	UPDATE CART_ITEM
	SET
		quantity = $3
	WHERE cart_id    = $1
	AND   product_id = $2
	RETURNING
		cart_id,
		product_id,
		quantity,
		unit_price
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		cartID,
		productID,
		quantity,
	)

	cartItem, err := scanCartItem(row)

	if err == pgx.ErrNoRows {
		cfg.Logs.Logger.Info(
			"Cart item not found for update",
			slog.Int("cart_id", cartID),
			slog.Int("product_id", productID),
			slog.String("func", "UpdateCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart_item{}, err
	}

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db update",
			slog.String("error", err.Error()),
			slog.String("func", "UpdateCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return model.Cart_item{}, fmt.Errorf("Error on db update: %w", err)
	}

	cfg.Logs.Logger.Info(fmt.Sprintf("Updated Cart Item for cart_id: %d, product_id: %d\n", cartItem.Cart_id, cartItem.Product_id))
	return cartItem, nil
}

func DeleteCartItem(cartID int, productID int) (int64, error) {
	cfg := config.GetConfig()

	sql := `
	DELETE FROM CART_ITEM
	WHERE cart_id   = $1
	AND   product_id = $2
	`

	result, err := cfg.DBConnection.Pool.Exec(
		cfg.DBConnection.Ctx,
		sql,
		cartID,
		productID,
	)

	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db delete",
			slog.String("error", err.Error()),
			slog.String("func", "DeleteCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return 0, fmt.Errorf("Error on db delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		cfg.Logs.Logger.Info(
			"Cart item not found for delete",
			slog.Int("cart_id", cartID),
			slog.Int("product_id", productID),
			slog.String("func", "DeleteCartItem"),
			slog.String("timestamp", time.Now().GoString()),
		)
		return 0, pgx.ErrNoRows
	}

	return result.RowsAffected(), nil
}
