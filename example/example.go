package main

import (
	"fmt"
	"github.com/dfryer1193/mjolnir/router"
	"github.com/dfryer1193/mjolnir/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	r := router.New()

	r.Get("/", utils.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) *utils.ApiError {
			_, err := w.Write([]byte("Hello World!"))
			if err != nil {
				return utils.InternalServerErr(err)
			}
			return nil
		}),
	)

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, r, 200, map[string]string{"msg": "Hello World!"})
	})

	r.Post("/json", utils.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) *utils.ApiError {
			var name struct {
				Name string `json:"name"`
			}

			_, err := utils.DecodeJSON(r, &name)
			if err != nil {
				return utils.BadRequestErr(err)
			}

			utils.RespondJSON(w, r, 200, map[string]string{"msg": "Hello " + name.Name})
			return nil
		}),
	)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("This is a panic")
	})

	r.Get("/error", utils.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) *utils.ApiError {
			return utils.NewApiError(fmt.Errorf("This is an error"), http.StatusServiceUnavailable)
		}),
	)

	log.Info().Msg("Server starting on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
	log.Info().Msg("Server stopped")
}
