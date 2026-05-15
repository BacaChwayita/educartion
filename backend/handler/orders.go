package handler

import (
	"net/http"
	"strconv"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func HandleGetOrderByID(w http.ResponseWriter, req *http.Request) {

	// TODO: This function currently does both `/api/orders/{id}` and `/api/orders/{id}/items` - which might be wrong? to think about and possibly split

	var id_string string
	var order_id int
	var err error

	id_string = req.PathValue("id")
	order_id, err = strconv.Atoi(id_string)
	if err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "Could not retrieve order id from path")
	}

	var response model.GetOrderByIDResponse

	order, err := db.GetOrderById(order_id)
	if err == pgx.ErrNoRows {
		rpc.WriteError(w, http.StatusNotFound, "No order with the specified id")
	}

	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving order")
	}

	response.Order_id = order.Order_id
	response.Order_number = order.Order_number
	response.Status = order.Order_number
	response.Subtotal_amount = order.Subtotal_amount
	response.Discount_amount = order.Discount_amount
	response.Total_amount = order.Total_amount
	response.Placed_at = order.Placed_at

	order_items, err := db.GetOrderItems(order.Order_id)
	if err == pgx.ErrNoRows {
		rpc.WriteError(w, http.StatusNotFound, "No order with the specified id")
	}

	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving order items")
	}

	for _, order_item := range order_items {
		product, err := db.GetProductById(order_item.Product_id)
		if err == pgx.ErrNoRows {
			rpc.WriteError(w, http.StatusNotFound, "No order with the specified id")
		}
		if err != nil {
			rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving product")
		}

		supplier, err := db.GetSupplierById(product.Supplier_id)
		if err == pgx.ErrNoRows {
			rpc.WriteError(w, http.StatusNotFound, "No order with the specified id")
		}
		if err != nil {
			rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving supplier")
		}

		item := model.Order_item_list{
			Quantity:        order_item.Quantity,
			Unit_price:      order_item.Unit_price,
			Discount_amount: order_item.Discount_amount,
			Product: model.Product_details{
				Product_id:       product.Product_id,
				Name:             product.Name,
				Description:      product.Description,
				Price:            product.Price,
				Discount_percent: product.Discount_percent,
				Stock_quantity:   product.Stock_quantity,
				Is_active:        product.Is_active,
				Supplier: model.Supplier_details{
					Supplier_id: supplier.Supplier_id,
					Name:        supplier.Name,
				},
			},
		}
		response.Order_item = append(response.Order_item, item)
	}

	rpc.WriteJSON(w, 200, response)
}
