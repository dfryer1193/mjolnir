package router

import (
	enhancedmiddleware "github.com/dfryer1193/mjolnir/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New creates a new pre-configured chi router
func New() *chi.Mux {
	r := chi.NewRouter()

	// Add default chi middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(enhancedmiddleware.RequestLogger)
	r.Use(enhancedmiddleware.ErrorHandler)

	return r
}
