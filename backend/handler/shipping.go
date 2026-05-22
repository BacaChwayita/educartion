package handler

import (
"errors"
"net/http"
"strconv"
"strings"

"github.com/P-SEN371-Group-3/educartion/rpc"
"github.com/P-SEN371-Group-3/educartion/service/db"
"github.com/jackc/pgx/v5"
)

func HandleShipmentByOrder(w http.ResponseWriter, req *http.Request) {
if req.Method != http.MethodGet {
http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
return
}

_, err := getBearerToken(req)
if err != nil {
rpc.WriteError(w, http.StatusUnauthorized, err.Error())
return
}

orderID := strings.TrimPrefix(req.URL.Path, "/api/shipments/")
if orderID == "" {
rpc.WriteError(w, http.StatusBadRequest, "order id is required")
return
}

id, err := strconv.Atoi(orderID)
if err != nil || id <= 0 {
rpc.WriteError(w, http.StatusBadRequest, "invalid order id")
return
}

shipment, err := db.GetShipmentByOrderId(id)
if err != nil {
if errors.Is(err, pgx.ErrNoRows) {
rpc.WriteError(w, http.StatusNotFound, "Shipment not found")
return
}
rpc.WriteError(w, http.StatusInternalServerError, "Error retrieving shipment")
return
}

rpc.WriteJSON(w, http.StatusOK, shipment)
}
