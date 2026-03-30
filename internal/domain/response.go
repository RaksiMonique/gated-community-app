package domain

import (
	"encoding/json"
	"net/http"
)

// Envelope is a generic structure for all API responses.
type Envelope struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Data    any      `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// WriteJSON is a standardized helper that writes a JSON response with optional headers.
func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	js = append(js, '\n')
	_, err = w.Write(js)
	return err
}
