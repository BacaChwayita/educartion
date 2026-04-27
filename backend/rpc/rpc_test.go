package rpc_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/P-SEN371-Group-3/educartion/rpc"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	rpc.WriteJSON(rr, 200, "All Ok!")

	if rr.Code != 200 {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var msg string
	err := json.NewDecoder(rr.Body).Decode(&msg)
	if err != nil {
		t.Errorf("Could not decode message")
	}
	if msg != "All Ok!" {
		t.Errorf("Expected 'All Ok!' got '%s'", msg)
	}

}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()

	rpc.WriteJSON(rr, 404, "Couldn't find something")

	if rr.Code != 404 {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}

	var msg string
	err := json.NewDecoder(rr.Body).Decode(&msg)
	if err != nil {
		t.Errorf("Could not decode message")
	}
	if msg != "Couldn't find something" {
		t.Errorf("Expected 'Couldn't find something' got '%s'", msg)
	}
}
