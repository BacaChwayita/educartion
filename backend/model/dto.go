package model

import "time"

// Note: Golang does not have a money type,
//       and due to floating point errors with float32 and float64,
//       we will be using int for the money values,
//       where 10.50 will be stored as 1050

type Account struct {
	Account_id     int       `json:"account_id"`
	Full_name      string    `json:"full_name"`
	Email          string    `json:"email"`
	Password_hash  string    `json:"password_hash"`
	Password_salt  string    `json:"password_salt"`
	Role           string    `json:"role"`
	Login_attempts int       `json:"login_attempts"`
	Is_active      bool      `json:"is_active"`
	Created_at     time.Time `json:"created_at"`
}

type AccountLogin struct {
	Token_id     int       `json:"token_id"`
	Account_id   int       `json:"account_id"`
	Token_string string    `json:"token_string"`
	Created_at   time.Time `json:"created_at"`
}

type Supplier struct {
	Supplier_id   int    `json:"supplier_id"`
	Name          string `json:"name"`
	Contact_email string `json:"contact_email"`
	Contact_phone string `json:"contact_phone"`
}

type Product struct {
	Product_id       int    `json:"product_id"`
	Supplier_id      int    `json:"supplier_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Price            int    `json:"price"`
	Discount_percent int    `json:"discount_percent"`
	Stock_quantity   int    `json:"stock_quantity"`
	Is_active        bool   `json:"is_active"`
}

type Product_image struct {
	Product_image_id int       `json:"product_image_id"`
	Product_id       int       `json:"product_id"`
	Image_url        string    `json:"image_url"`
	Alt_text         string    `json:"alt_text"`
	Is_primary       bool      `json:"is_primary"`
	Sort_order       int       `json:"sort_order"`
	Created_at       time.Time `json:"created_at"`
}

type Cart struct {
	Cart_id     int       `json:"cart_id"`
	Account_id  int       `json:"account_id"`
	Session_key string    `json:"session_key"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type Orders struct {
	Order_id        int       `json:"order_id"`
	Account_id      int       `json:"account_id"`
	Order_number    string    `json:"order_number"`
	Status          string    `json:"status"`
	Subtotal_amount int       `json:"subtotal_amount"`
	Discount_amount int       `json:"discount_amount"`
	Total_amount    int       `json:"total_amount"`
	Placed_at       time.Time `json:"placed_at"`
}

type Payment struct {
	Payment_id      int       `json:"payment_id"`
	Order_id        int       `json:"order_id"`
	Payment_method  string    `json:"payment_method"`
	Payment_status  string    `json:"payment_status"`
	Transaction_ref string    `json:"transaction_ref"`
	Amount          int       `json:"amount"`
	Paid_at         time.Time `json:"paid_at"`
}

type Shipment struct {
	Shipment_id            int       `json:"shipment_id"`
	Order_id               int       `json:"order_id"`
	Courier_name           string    `json:"courier_name"`
	Courier_type           string    `json:"courier_type"`
	Delivery_reference     string    `json:"delivery_reference"`
	Delivery_address_line1 string    `json:"delivery_address_line1"`
	Delivery_address_line2 string    `json:"delivery_address_line2"`
	Delivery_city          string    `json:"delivery_city"`
	Delivery_state         string    `json:"delivery_state"`
	Delivery_postal_code   string    `json:"delivery_postal_code"`
	Delivery_country       string    `json:"delivery_country"`
	Recipient_name         string    `json:"recipient_name"`
	Recipient_phone        string    `json:"recipient_phone"`
	Shipment_status        string    `json:"shipment_status"`
	Dispatched_at          time.Time `json:"dispatched_at"`
	Delivered_at           time.Time `json:"delivered_at"`
	Created_at             time.Time `json:"created_at"`
}

type Cart_item struct {
	Cart_id    int `json:"cart_id"`
	Product_id int `json:"product_id"`
	Quantity   int `json:"quantity"`
	Unit_price int `json:"unit_price"`
}

type Order_item struct {
	Order_id        int `json:"order_id"`
	Product_id      int `json:"product_id"`
	Quantity        int `json:"quantity"`
	Unit_price      int `json:"unit_price"`
	Discount_amount int `json:"discount_amount"`
}
