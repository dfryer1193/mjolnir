package middleware

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestErrorHandler(t *testing.T) {
    tests := []struct {
        name           string
        handler       func(w http.ResponseWriter, r *http.Request)
        expectedCode  int
        expectedError string
        checkLogs     bool
        expectedLogs  []string
    }{
        {
            name: "no error",
            handler: func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            },
            expectedCode:  http.StatusOK,
            expectedError: "",
        },
        {
            name: "bad request error",
            handler: func(w http.ResponseWriter, r *http.Request) {
                SetBadRequestError(r, errors.New("invalid input"))
            },
            expectedCode:  http.StatusBadRequest,
            expectedError: "invalid input",
        },
        {
            name: "internal server error",
            handler: func(w http.ResponseWriter, r *http.Request) {
                SetInternalError(r, errors.New("database connection failed"))
            },
            expectedCode:  http.StatusInternalServerError,
            expectedError: "Internal Server Error",
            checkLogs:     true,
            expectedLogs:  []string{"database connection failed", "internal server error occurred"},
        },
        {
            name: "not found error",
            handler: func(w http.ResponseWriter, r *http.Request) {
                SetNotFoundError(r, errors.New("resource not found"))
            },
            expectedCode:  http.StatusNotFound,
            expectedError: "resource not found",
        },
        {
            name: "unauthorized error",
            handler: func(w http.ResponseWriter, r *http.Request) {
                SetUnauthorizedError(r, errors.New("invalid token"))
            },
            expectedCode:  http.StatusUnauthorized,
            expectedError: "invalid token",
        },
        {
            name: "panic recovery",
            handler: func(w http.ResponseWriter, r *http.Request) {
                panic("unexpected panic")
            },
            expectedCode:  http.StatusInternalServerError,
            expectedError: "Internal Server Error",
            checkLogs:     true,
            expectedLogs:  []string{"panic: unexpected panic"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Capture logs if needed
            var buf bytes.Buffer
            if tt.checkLogs {
                log.Logger = zerolog.New(&buf)
            }

            // Create test handler
            handler := ErrorHandler(http.HandlerFunc(tt.handler))

            // Create test request with request ID
            req := httptest.NewRequest(http.MethodGet, "/test", nil)
            ctx := context.WithValue(req.Context(), requestIDKey, "test-request-id")
            req = req.WithContext(ctx)

            // Create response recorder
            rr := httptest.NewRecorder()

            // Execute request
            handler.ServeHTTP(rr, req)

            // Check status code
            if rr.Code != tt.expectedCode {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    rr.Code, tt.expectedCode)
            }

            // Check response
            if tt.expectedError != "" {
                var errResp ErrorResponse
                if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
                    t.Fatalf("failed to decode response: %v", err)
                }

                if errResp.Error != tt.expectedError {
                    t.Errorf("unexpected error message: got %q want %q",
                        errResp.Error, tt.expectedError)
                }

                if errResp.Code != tt.expectedCode {
                    t.Errorf("unexpected error code: got %d want %d",
                        errResp.Code, tt.expectedCode)
                }
            }

            // Check Content-Type header
            if tt.expectedError != "" {
                contentType := rr.Header().Get("Content-Type")
                if contentType != "application/json" {
                    t.Errorf("wrong Content-Type: got %v want application/json", contentType)
                }
            }

            // Check logs
            if tt.checkLogs {
                logStr := buf.String()
                for _, expectedLog := range tt.expectedLogs {
                    if !strings.Contains(logStr, expectedLog) {
                        t.Errorf("log doesn't contain %q\nLog: %s", expectedLog, logStr)
                    }
                }
            }
        })
    }
}

func TestSetError(t *testing.T) {
    tests := []struct {
        name       string
        status     int
        err        error
        isInternal bool
    }{
        {
            name:       "internal error",
            status:     http.StatusInternalServerError,
            err:        errors.New("internal error"),
            isInternal: true,
        },
        {
            name:       "client error",
            status:     http.StatusBadRequest,
            err:        errors.New("bad request"),
            isInternal: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, "/test", nil)
            SetError(req, tt.status, tt.err)

            reqErr, ok := req.Context().Value(errorCtxKey).(*RequestError)
            if !ok {
                t.Fatal("error not set in context")
            }

            if reqErr.Status != tt.status {
                t.Errorf("wrong status: got %d want %d", reqErr.Status, tt.status)
            }

            if reqErr.Err.Error() != tt.err.Error() {
                t.Errorf("wrong error: got %v want %v", reqErr.Err, tt.err)
            }

            if reqErr.IsInternal != tt.isInternal {
                t.Errorf("wrong IsInternal flag: got %v want %v", 
                    reqErr.IsInternal, tt.isInternal)
            }
        })
    }
}

func BenchmarkErrorHandler(b *testing.B) {
    handler := ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        SetError(r, http.StatusBadRequest, errors.New("test error"))
    }))
    
    req := httptest.NewRequest(http.MethodGet, "/benchmark", nil)
    rr := httptest.NewRecorder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        handler.ServeHTTP(rr, req)
    }
}