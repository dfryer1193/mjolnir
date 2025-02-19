package utils

import (
	"encoding/json"
	"errors"
	"github.com/dfryer1193/mjolnir/middleware"
	"net/http"
)

// RespondJSON sends a JSON response with proper headers
func RespondJSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		middleware.SetInternalError(r, err)
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
