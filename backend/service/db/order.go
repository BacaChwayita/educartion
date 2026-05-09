package model

import "time"

type Order struct {
	Order_id    int
	Account_id  int
	Supplier_id int
	Status      string
	Notes       string
	Created_at  time.Time
	Updated_at  time.Time
}

type OrderItem struct {
	Order_item_id int
	Order_id      int
	Product_name  string
	Quantity      int
	Unit_price    float64
}
