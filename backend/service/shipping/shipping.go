package shipping

import (
	"errors"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func GetShipmentByOrderID(orderID int) (model.Shipment, error) {
	shipment, err := db.GetShipmentByOrderId(orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Shipment{}, ErrShipmentNotFound
		}
		return model.Shipment{}, err
	}

	return shipment, nil
}
