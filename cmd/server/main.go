package main

import (
	"net/http"

	"github.com/anatolev-max/metrics-tpl/cmd/common"
	"github.com/anatolev-max/metrics-tpl/cmd/handlers"
	"github.com/anatolev-max/metrics-tpl/cmd/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	memStorage := storage.NewMemStorage()

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", handlers.GetMainWebhook(memStorage))
		router.Get(common.GetEndpoint+"/{type}/{name}", handlers.GetValueWebhook(memStorage))
		router.Post(common.UpdateEndpoint+"/{type}/{name}/{value}", handlers.GetUpdateWebhook(memStorage))
	})

	return http.ListenAndServe(common.ServerPort, router)
}
