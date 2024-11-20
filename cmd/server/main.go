package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/anatolev-max/metrics-tpl/config"
	"github.com/anatolev-max/metrics-tpl/internal/enum"
	"github.com/anatolev-max/metrics-tpl/internal/handlers"
	"github.com/anatolev-max/metrics-tpl/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	c := config.NewConfig()
	parseFlags(c)

	if err := run(c); err != nil {
		panic(err)
	}
}

func run(c config.Config) error {
	s := storage.NewMemStorage()
	h := handlers.NewHTTPHandler(s)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", h.GetMainWebhook())
		router.Get(enum.GetEndpoint.String()+"/{type}/{name}", h.GetValueWebhook())
		router.Post(enum.UpdateEndpoint.String()+"/{type}/{name}/{value}", h.GetUpdateWebhook(c))
	})

	fmt.Println("Running server on", options.runAddr)
	u, _ := url.ParseRequestURI(options.runAddr)

	return http.ListenAndServe(":"+u.Port(), router)
}
