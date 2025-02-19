package utils

import (
	"encoding/json"
	"errors"
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
func DecodeJSON(r *http.Request, v interface{}) {
	if ValidateContentType(r, "application/json") {
		middleware.SetBadRequestError(r, errors.New("Content-Type must be application/json"))
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		middleware.SetBadRequestError(r, err)
	}
}

// ValidateContentType checks if the request has the required content type
func ValidateContentType(r *http.Request, contentType string) bool {
	return r.Header.Get("Content-Type") == contentType
}
