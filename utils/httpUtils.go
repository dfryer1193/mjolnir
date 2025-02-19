package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dfryer1193/mjolnir/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

// RespondJSON sends a JSON response with proper headers
func RespondJSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		middleware.SetInternalError(r, err)
		return
	}

	w.WriteHeader(status)

	if _, err := w.Write(response); err != nil {
		log.Error().Err(fmt.Errorf("failed to write response: %w", err)).Msg("Failed to write response")
	}
}

// DecodeJSON decodes JSON request body into the provided struct
func DecodeJSON(r *http.Request, v interface{}) error {
	if !ValidateContentType(r, "application/json") {
		return fmt.Errorf("Content-Type %s is not supported", r.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// ValidateContentType checks if the request has the required content type
func ValidateContentType(r *http.Request, contentType string) bool {
	return r.Header.Get("Content-Type") == contentType
}
