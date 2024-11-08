package main

import (
	"net/http"

	"github.com/anatolev-max/metrics-tpl/internal/config"
	"github.com/anatolev-max/metrics-tpl/internal/handlers"
	"github.com/anatolev-max/metrics-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := storage.NewMemStorage()

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", handlers.GetMainWebhook(s))
		router.Get(config.GetEndpoint+"/{type}/{name}", handlers.GetValueWebhook(s))
		router.Post(config.UpdateEndpoint+"/{type}/{name}/{value}", handlers.GetUpdateWebhook(s))
	})

	return http.ListenAndServe(config.ServerPort, router)
}
