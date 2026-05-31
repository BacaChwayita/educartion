package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/jackc/pgx/v5"
)

const productImageColumns = `
	product_image_id,
	product_id,
	image_url,
	alt_text,
	is_primary,
	sort_order,
	created_at
`

func GetProductImagesByProductId(product_id int) ([]model.Product_image, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + productImageColumns +
		`
		FROM product_image
		WHERE product_id = $1
		ORDER BY is_primary DESC, sort_order ASC, product_image_id ASC
	`

	rows, err := cfg.DBConnection.Pool.Query(
		cfg.DBConnection.Ctx,
		sql,
		product_id,
	)
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetProductImagesByProductId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", product_id),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}
	defer rows.Close()

	images := make([]model.Product_image, 0)
	for rows.Next() {
		var image model.Product_image
		err := rows.Scan(
			&image.Product_image_id,
			&image.Product_id,
			&image.Image_url,
			&image.Alt_text,
			&image.Is_primary,
			&image.Sort_order,
			&image.Created_at,
		)
		if err != nil {
			cfg.Logs.Logger.Info(
				"Error on db select",
				slog.String("error", err.Error()),
				slog.String("func", "GetProductImagesByProductId"),
				slog.String("timestamp", time.Now().GoString()),
				slog.Int("product_id", product_id),
			)
			return nil, fmt.Errorf("Error on db select: %w", err)
		}

		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetProductImagesByProductId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", product_id),
		)
		return nil, fmt.Errorf("Error on db select: %w", err)
	}

	return images, nil
}

func GetPrimaryProductImageByProductId(product_id int) (model.Product_image, error) {
	cfg := config.GetConfig()

	sql := `SELECT ` + productImageColumns +
		`
		FROM product_image
		WHERE product_id = $1
		ORDER BY is_primary DESC, sort_order ASC, product_image_id ASC
		LIMIT 1
	`

	row := cfg.DBConnection.Pool.QueryRow(
		cfg.DBConnection.Ctx,
		sql,
		product_id,
	)

	var image model.Product_image
	err := row.Scan(
		&image.Product_image_id,
		&image.Product_id,
		&image.Image_url,
		&image.Alt_text,
		&image.Is_primary,
		&image.Sort_order,
		&image.Created_at,
	)
	if err == pgx.ErrNoRows {
		return model.Product_image{}, err
	}
	if err != nil {
		cfg.Logs.Logger.Info(
			"Error on db select",
			slog.String("error", err.Error()),
			slog.String("func", "GetPrimaryProductImageByProductId"),
			slog.String("timestamp", time.Now().GoString()),
			slog.Int("product_id", product_id),
		)
		return model.Product_image{}, fmt.Errorf("Error on db select: %w", err)
	}

	return image, nil
}
