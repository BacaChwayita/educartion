package db_test

import (
	"testing"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

// TODO: These tests are currently failing due to Order_id not existing in the orders table, breaking the fk constraint
//       Need to refactor to InsertOrder, which also needs an account_id
//       Alternatively need dummy data in test db to refer to, ie order_id 1 and account_id 1 always existing there. (setup_db scripts)

func TestInsertShipment(t *testing.T) {
	newShipment := model.Shipment{
		Shipment_id:            -1,
		Order_id:               -1,
		Courier_name:           "test name",
		Courier_type:           "test type",
		Delivery_reference:     "some ref",
		Delivery_address_line1: "10 bard street",
		Delivery_address_line2: "Bards Alley",
		Delivery_city:          "Bards City",
		Delivery_state:         "BardState",
		Delivery_postal_code:   "1234",
		Delivery_country:       "Bards Country",
		Recipient_name:         "John Doe",
		Recipient_phone:        "+27 12 345 6789",
		Shipment_status:        "pending",
		Dispatched_at:          time.Now(),
		Delivered_at:           time.Now(),
		Created_at:             time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	ship, err := db.InsertShipment(newShipment)

	if err != nil {
		t.Errorf("InsertShipment() call failed: %s", err.Error())
	}

	if ship.Shipment_id <= 0 {
		t.Errorf("Expected valid Shipment_id, Got: %d", ship.Shipment_id)
	}
}

func TestGetShiplierById(t *testing.T) {
	insShip := model.Shipment{
		Shipment_id:            -1,
		Order_id:               -1,
		Courier_name:           "test name",
		Courier_type:           "test type",
		Delivery_reference:     "some ref",
		Delivery_address_line1: "10 bard street",
		Delivery_address_line2: "Bards Alley",
		Delivery_city:          "Bards City",
		Delivery_state:         "BardState",
		Delivery_postal_code:   "1234",
		Delivery_country:       "Bards Country",
		Recipient_name:         "John Doe",
		Recipient_phone:        "+27 12 345 6789",
		Shipment_status:        "pending",
		Dispatched_at:          time.Now(),
		Delivered_at:           time.Now(),
		Created_at:             time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newShip, err := db.InsertShipment(insShip)
	if err != nil {
		t.Fatalf("InsertShipment() call failed: %s", err.Error())
	}

	getShip, err := db.GetShipmentById(newShip.Shipment_id)

	if err != nil {
		t.Errorf("Error not nil, got %s", err.Error())
	}

	if getShip.Shipment_id != newShip.Shipment_id {
		t.Errorf(
			"Got %d, Expected %d. Expected GetShipment to return same shipment_id as InsertShipment inserted.",
			getShip.Shipment_id,
			newShip.Shipment_id,
		)
	}

}

func TestUpdateShipment(t *testing.T) {
	insShip := model.Shipment{
		Shipment_id:            -1,
		Order_id:               -1,
		Courier_name:           "test name",
		Courier_type:           "test type",
		Delivery_reference:     "some ref",
		Delivery_address_line1: "10 bard street",
		Delivery_address_line2: "Bards Alley",
		Delivery_city:          "Bards City",
		Delivery_state:         "BardState",
		Delivery_postal_code:   "1234",
		Delivery_country:       "Bards Country",
		Recipient_name:         "John Doe",
		Recipient_phone:        "+27 12 345 6789",
		Shipment_status:        "pending",
		Dispatched_at:          time.Now(),
		Delivered_at:           time.Now(),
		Created_at:             time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newShip, err := db.InsertShipment(insShip)
	if err != nil {
		t.Fatalf("InsertShipment() call failed: %s", err.Error())
	}

	chgShip := newShip
	chgShip.Shipment_status = "delivered"
	rowCount, err := db.UpdateShipment(chgShip)
	if err != nil {
		t.Errorf("UpdateShipment() call failed: %s", err.Error())
	}

	updatedShip, err := db.GetShipmentById(chgShip.Shipment_id)

	if rowCount != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowCount)
	}

	if updatedShip.Shipment_status != "delivered" {
		t.Errorf("Expected 'delivered', Got %s", updatedShip.Shipment_status)
	}
}

func TestDeleteShipmentById(t *testing.T) {
	insShip := model.Shipment{
		Shipment_id:            -1,
		Order_id:               -1,
		Courier_name:           "test name",
		Courier_type:           "test type",
		Delivery_reference:     "some ref",
		Delivery_address_line1: "10 bard street",
		Delivery_address_line2: "Bards Alley",
		Delivery_city:          "Bards City",
		Delivery_state:         "BardState",
		Delivery_postal_code:   "1234",
		Delivery_country:       "Bards Country",
		Recipient_name:         "John Doe",
		Recipient_phone:        "+27 12 345 6789",
		Shipment_status:        "pending",
		Dispatched_at:          time.Now(),
		Delivered_at:           time.Now(),
		Created_at:             time.Now(),
	}

	err := SetupTestDBConfigAndConnection()
	if err != nil {
		t.Fatalf("SetupTestDBConfigAndConnection() call failed: %s", err.Error())
	}

	newShip, err := db.InsertShipment(insShip)
	if err != nil {
		t.Fatalf("InsertShipment() call failed: %s", err.Error())
	}

	rowCount, err := db.DeleteShipmentById(newShip.Shipment_id)
	if err != nil {
		t.Errorf("DeleteShipmentById() call failed: %s", err.Error())
	}

	if rowCount != 1 {
		t.Errorf("Expected 1, got %d", rowCount)
	}

	deletedShip, err := db.GetShipmentById(newShip.Shipment_id)
	if err != pgx.ErrNoRows {
		t.Errorf("Expected pgx.ErrNoRows error to be returned, Got: %s", err.Error())
	}

	if deletedShip.Shipment_id != 0 {
		t.Errorf("Expected 0, Got %d", deletedShip.Shipment_id)
	}
}
