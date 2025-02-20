package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondJSON sends a JSON response with proper headers
func RespondJSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	w.WriteHeader(status)

	if _, err := w.Write(response); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

// DecodeJSON decodes JSON request body into the provided struct
func DecodeJSON(r *http.Request, v interface{}) error {
	if !ValidateContentType(r, "application/json") {
		return fmt.Errorf("Content-Type %s is not supported", r.Header.Get("Content-Type"))
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// ValidateContentType checks if the request has the required content type
func ValidateContentType(r *http.Request, contentType string) bool {
	return r.Header.Get("Content-Type") == contentType
}
