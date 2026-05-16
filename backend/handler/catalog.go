import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	"github.com/P-SEN371-Group-3/educartion/service/catalog"
)
//When a user clicks on catalog but have not searched 
func HandleGetProducts(w http.ResponseWriter, req *http.Request) {

	var pr model.ProductRequest
	err := json.NewDecoder(req.Body).Decode(&pr)
    if err != nil{
    rpc.WriteError(
			w,
			http.StatusBadRequest,
			rpc.ErrDecodeHTTPRequestBody,
		)
    return

	}

	products, err := catalog.GetProducts()

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve products",
		)

		return
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

	product, err := catalog.GetProductById(productID)

	if err != nil {

		rpc.WriteError(
			w,
			http.StatusNotFound,
			"Product not found",
		)

		return
	}

	rpc.WriteJSON(
		w,
		http.StatusOK,
		product,
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

	product, err := catalog.CreateProduct(newProduct)

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

	rowCount, err := catalog.UpdateProduct(updatedProduct)

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

	rowCount, err := catalog.DeleteProduct(productID)

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
