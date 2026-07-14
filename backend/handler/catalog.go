package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

// When a user clicks on catalog but have not searched.
func HandleGetProducts(w http.ResponseWriter, req *http.Request) {

	var pr model.ProductRequest
	err := json.NewDecoder(req.Body).Decode(&pr)
	if err != nil {
		rpc.WriteError(
			w,
			http.StatusBadRequest,
			rpc.ErrDecodeHTTPRequestBody,
		)
		return

	}

	products, err := db.GetAllProducts()

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve products",
		)

		return
	}

	for index, product := range products {
		image, err := db.GetPrimaryProductImageByProductId(product.Product_id)
		if err == pgx.ErrNoRows {
			continue
		}
		if err != nil {
			rpc.WriteError(
				w,
				http.StatusInternalServerError,
				"Failed to retrieve product images",
			)
			return
		}

		products[index].Image_url = image.Image_url
	}

	rpc.WriteJSON(
		w,
		http.StatusOK,
		products,
	)
}
func HandleGetProductById(w http.ResponseWriter, req *http.Request) {

	productIDString := req.URL.Query().Get("id")

	productID, err := strconv.Atoi(productIDString)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid product id",
		)

		return
	}

	product, err := db.GetProductById(productID)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusNotFound,
			"Product not found",
		)

		return
	}

	supplier, err := db.GetSupplierById(product.Supplier_id)
	if err == pgx.ErrNoRows {
		rpc.WriteError(
			w,
			http.StatusNotFound,
			"Supplier not found",
		)
		return
	}
	if err != nil {
		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve supplier",
		)
		return
	}

	productImages, err := db.GetProductImagesByProductId(product.Product_id)
	if err != nil {
		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve product images",
		)
		return
	}

	response := model.Product_details{
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
		Product_images: productImages,
	}

	rpc.WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}
func HandleCreateProduct(w http.ResponseWriter, req *http.Request) {

	var newProduct model.Product

	err := json.NewDecoder(req.Body).Decode(&newProduct)

	if err != nil {

		logErrDecodeBody(
			"HandleCreateProduct",
			req,
			err,
		)

		rpc.WriteError(
			w,
			http.StatusBadRequest,
			rpc.ErrDecodeHTTPRequestBody,
		)

		return
	}

	product, err := db.InsertProduct(newProduct)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	rpc.WriteJSON(
		w,
		http.StatusCreated,
		product,
	)
}
func HandleUpdateProduct(w http.ResponseWriter, req *http.Request) {

	productIDString := req.URL.Query().Get("id")

	productID, err := strconv.Atoi(productIDString)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid product id",
		)

		return
	}

	var updatedProduct model.Product

	err = json.NewDecoder(req.Body).Decode(&updatedProduct)

	if err != nil {

		logErrDecodeBody(
			"HandleUpdateProduct",
			req,
			err,
		)

		rpc.WriteError(
			w,
			http.StatusBadRequest,
			rpc.ErrDecodeHTTPRequestBody,
		)

		return
	}

	updatedProduct.Product_id = productID

	rowCount, err := db.UpdateProduct(updatedProduct)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	rpc.WriteJSON(
		w,
		http.StatusOK,
		map[string]any{
			"rows_affected": rowCount,
		},
	)
}
func HandleDeleteProduct(w http.ResponseWriter, req *http.Request) {

	productIDString := req.URL.Query().Get("id")

	productID, err := strconv.Atoi(productIDString)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid product id",
		)

		return
	}

	rowCount, err := db.DeleteProductById(productID)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	rpc.WriteJSON(
		w,
		http.StatusOK,
		map[string]any{
			"rows_affected": rowCount,
		},
	)
}
