package main

import (
	"fmt"
	"net/http"
	. "net/url"

	. "github.com/anatolev-max/metrics-tpl/config"
	"github.com/anatolev-max/metrics-tpl/internal/enum"
	"github.com/anatolev-max/metrics-tpl/internal/handlers"
	"github.com/anatolev-max/metrics-tpl/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	c := NewConfig()
	parseFlags(c)

	if err := run(c); err != nil {
		panic(err)
	}
}

func run(c Config) error {
	s := storage.NewMemStorage()

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", handlers.GetMainWebhook(s))
		router.Get(enum.GetEndpoint.String()+"/{type}/{name}", handlers.GetValueWebhook(s))
		router.Post(enum.UpdateEndpoint.String()+"/{type}/{name}/{value}", handlers.GetUpdateWebhook(s, c))
	})

	fmt.Println("Running server on", options.flagRunAddr)
	url, _ := ParseRequestURI(options.flagRunAddr)

	return http.ListenAndServe(":"+url.Port(), router)
}
