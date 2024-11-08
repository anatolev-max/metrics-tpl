package handlers

import (
	"encoding/json"
	"github.com/anatolev-max/metrics-tpl/internal/render"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/anatolev-max/metrics-tpl/internal/config"
	"github.com/anatolev-max/metrics-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func GetMainWebhook(s storage.MemStorage) func(http.ResponseWriter, *http.Request) {
	return render.IncludeTemplate("index.html", map[string]any{
		"Title":    "Metrics-tpl",
		"Counters": s.Counter,
		"Gauges":   s.Gauge,
	})
}

func GetValueWebhook(s storage.MemStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var data []byte
		metricType := strings.ToLower(chi.URLParam(req, "type"))
		metricName := strings.ToLower(chi.URLParam(req, "name"))

		supportedMTypes := []string{config.Counter, config.Gauge}
		if !slices.Contains(supportedMTypes, metricType) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		msVal := reflect.ValueOf(s)

		for fieldIndex := 0; fieldIndex < msVal.NumField(); fieldIndex++ {
			if field := msVal.Type().Field(fieldIndex).Name; strings.ToLower(field) == metricType {
				iter := msVal.FieldByName(field).MapRange()
				for iter.Next() {
					if strings.ToLower(iter.Key().String()) == metricName {
						data, _ = json.Marshal(iter.Value().Interface())
					}
				}
			}
		}

		if len(data) == 0 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", config.TextPlain)
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(data); err != nil {
			panic(err)
		}
	}
}

func GetUpdateWebhook(ms storage.MemStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// TODO: chi.URLParam
		urlPath := strings.TrimLeft(req.RequestURI, config.ServerHost+config.ServerPort)
		urlParams := strings.Split(urlPath, "/")
		if len(urlParams) != 4 {
			return
		}

		metricType := urlParams[1]
		metricName := urlParams[2]
		metricValue := urlParams[3]

		if metricName == "" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		supportedMTypes := []string{config.Counter, config.Gauge}
		if req.Header.Get("Content-Type") != config.TextPlain || !slices.Contains(supportedMTypes, metricType) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var convMetricValue interface{}
		var err error

		switch metricType {
		case config.Counter:
			convMetricValue, err = strconv.ParseInt(metricValue, 0, 64)
		case config.Gauge:
			convMetricValue, err = strconv.ParseFloat(metricValue, 64)
		}

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", config.TextPlain)
		res.WriteHeader(http.StatusOK)
		ms.UpdateServerData(metricName, convMetricValue)
	}
}
