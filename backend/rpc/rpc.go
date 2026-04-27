package rpc

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		// TODO: log?
		return err
	}

	return nil
}

func WriteError(w http.ResponseWriter, status int, msg string) error {
	return WriteJSON(w, status, map[string]string{"error": msg})
}
