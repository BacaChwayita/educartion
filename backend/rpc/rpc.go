package rpc

import (
	"encoding/json"
	"net/http"
)

// Possible errors that will be written through WriteError
const (
	ErrDecodeHTTPRequestBody string = "Failed to decode HTTP request body"
	ErrNoAuthorizationKey    string = "No valid authorization key found"
)

// WriteJSON writes a response back to the requester.
//
// @w -> the writer to write the message to
//
// @status -> the http status code to write back (ie 200 or http.StatusOK)
//
// @msg -> nil or the ResponseType to send back. Ensure that type has json tags,
// or pass through a value that contains the mapping.
func WriteJSON(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if msg != nil {
		err := json.NewEncoder(w).Encode(msg)
		if err != nil {
			// TODO: log? not just error but that message sent...
			return err
		}
	}

	return nil
}

// WriteError writes an error response back to the requester, using WriteJSON.
//
// @w -> the writer to write the message to
//
// @status -> the http status code to write back (ie 200 or http.StatusOK)
//
// @msg -> the error message to write, will be attached to json tag "error"
func WriteError(w http.ResponseWriter, status int, msg string) error {
	return WriteJSON(w, status, map[string]string{"error": msg})
}
