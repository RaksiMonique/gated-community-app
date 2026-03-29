package handlers

import (
	"encoding/json"
	"gated-community-api/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := HandleHealth(nil) // In a real test, you might pass a mock or actual test DB

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify Content-Type header
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Decode the JSON to verify structure and data
	var response domain.Envelope
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("expected success to be true, got %v", response.Success)
	}

	// Type assertion for the Data field
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data field is not a map: %T", response.Data)
	}

	if dataMap["status"] != "ok" || dataMap["database"] != "disconnected" {
		t.Errorf("handler returned unexpected data for nil DB: %+v", response.Data)
	}
}
