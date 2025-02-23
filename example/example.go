package main

import (
	"errors"
	"github.com/dfryer1193/mjolnir/middleware"
	"github.com/dfryer1193/mjolnir/router"
	"github.com/dfryer1193/mjolnir/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	r := router.New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to write response")
		}
	})

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		err := utils.RespondJSON(w, r, 200, map[string]string{"msg": "Hello World!"})
		if err != nil {
			middleware.SetInternalError(r, err)
		}
	})

	r.Post("/json", func(w http.ResponseWriter, r *http.Request) {
		var name struct {
			Name string `json:"Name"`
		}

		err := utils.DecodeJSON(r, &name)
		if err != nil {
			middleware.SetBadRequestError(r, err)
			return
		}

		err = utils.RespondJSON(w, r, 200, map[string]string{"msg": "Hello " + name.Name})
		if err != nil {
			middleware.SetInternalError(r, err)
		}
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("This is a panic")
	})

	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		middleware.SetError(r, 504, errors.New("this is an error"))
	})

	log.Info().Msg("Server starting on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
	log.Info().Msg("Server stopped")
}
