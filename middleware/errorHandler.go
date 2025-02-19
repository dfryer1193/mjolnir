package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// RequestError holds the error information for a request
type RequestError struct {
	Status     int
	Err        error
	IsInternal bool
}

// ErrorResponse represents the standard error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// SetError stores an error in the request context
func SetError(r *http.Request, status int, err error) {
	reqErr := &RequestError{
		Status:     status,
		Err:        err,
		IsInternal: status >= http.StatusInternalServerError,
	}
	*r = *r.WithContext(context.WithValue(r.Context(), errorCtxKey, reqErr))
}

func SetInternalError(r *http.Request, err error) {
	SetError(r, http.StatusInternalServerError, err)
}

func SetBadRequestError(r *http.Request, err error) {
	SetError(r, http.StatusBadRequest, err)
}

func SetNotFoundError(r *http.Request, err error) {
	SetError(r, http.StatusNotFound, err)
}

func SetUnauthorizedError(r *http.Request, err error) {
	SetError(r, http.StatusUnauthorized, err)
}

// ErrorHandler middleware processes any errors set during request handling
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Handle panics as internal errors
				SetError(r, http.StatusInternalServerError, fmt.Errorf("panic: %v", rec))
			}

			// Check if an error was set in the request context
			if err, ok := r.Context().Value(errorCtxKey).(*RequestError); ok {
				handleError(w, r, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, r *http.Request, reqErr *RequestError) {
	if reqErr.IsInternal {
		// Log internal errors with full details
		log.Error().
			Str("request_id", GetRequestID(r.Context())).
			Err(reqErr.Err).
			Int("status", reqErr.Status).
			Str("path", r.URL.Path).
			Str("method", r.Method).
			Msg("internal server error occurred")

		// Return generic error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Internal Server Error",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	// For non-internal errors, return the actual error
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reqErr.Status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: reqErr.Err.Error(),
		Code:  reqErr.Status,
	})
}
