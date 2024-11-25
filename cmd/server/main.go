package main

import (
	"net/http"
	"net/url"

	"github.com/anatolev-max/metrics-tpl/config"
	"github.com/anatolev-max/metrics-tpl/internal/enum"
	"github.com/anatolev-max/metrics-tpl/internal/handlers"
	"github.com/anatolev-max/metrics-tpl/internal/logger"
	"github.com/anatolev-max/metrics-tpl/internal/storage"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	c := config.NewConfig()
	parseFlags(c)

	if err := run(c); err != nil {
		panic(err)
	}
}

func run(c config.Config) error {
	if err := logger.Initialize(c.LogLevel); err != nil {
		return err
	}

	s := storage.NewMemStorage()
	h := handlers.NewHTTPHandler(s)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get(
			"/",
			logger.RequestLogger(h.GetMainWebhook()),
		)
		router.Get(
			enum.GetEndpoint.String()+"/{type}/{name}",
			logger.RequestLogger(h.GetValueWebhook()),
		)
		router.Post(
			enum.UpdateEndpoint.String()+"/{type}/{name}/{value}",
			logger.RequestLogger(h.GetUpdateWebhook(c)),
		)
	})

	u, _ := url.ParseRequestURI(options.runAddr)
	logger.Log.Info("Running server", zap.String("address", options.runAddr))

	return http.ListenAndServe(":"+u.Port(), router)
}
