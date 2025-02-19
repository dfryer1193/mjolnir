package main

import (
	"github.com/dfryer1193/mjolnir/router"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	r := router.New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("This is a panic")
	})

	log.Info().Msg("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
