package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/anatolev-max/metrics-tpl/config"
	"github.com/anatolev-max/metrics-tpl/internal/enum"
	"github.com/anatolev-max/metrics-tpl/internal/render"
	"github.com/anatolev-max/metrics-tpl/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Webhook interface {
	GetMainWebhook() func(http.ResponseWriter, *http.Request)
	GetValueWebhook() func(http.ResponseWriter, *http.Request)
	GetUpdateWebhook(c config.Config) func(http.ResponseWriter, *http.Request)
}

type HTTPHandler struct {
	memStorage *storage.MemStorage
}

func NewHTTPHandler(s *storage.MemStorage) *HTTPHandler {
	return &HTTPHandler{
		memStorage: s,
	}
}

func (h *HTTPHandler) GetMainWebhook() func(http.ResponseWriter, *http.Request) {
	return render.IncludeTemplate("index.html", map[string]any{
		"Title":    "Metrics-tpl",
		"Counters": h.memStorage.Counter,
		"Gauges":   h.memStorage.Gauge,
	})
}

func (h *HTTPHandler) GetValueWebhook() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var data []byte
		metricType := strings.ToLower(chi.URLParam(req, "type"))
		metricName := strings.ToLower(chi.URLParam(req, "name"))

		supportedMTypes := []string{enum.Counter.String(), enum.Gauge.String()}
		if !slices.Contains(supportedMTypes, metricType) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		msVal := reflect.ValueOf(h.memStorage).Elem()

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

		res.Header().Set("Content-Type", enum.TextPlain.String())
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(data); err != nil {
			panic(err)
		}
	}
}

func (h *HTTPHandler) GetUpdateWebhook(c config.Config) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// TODO: chi.URLParam
		urlPath := strings.TrimLeft(req.RequestURI, c.Server.Scheme+c.Server.Host+c.Server.Port)
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

		supportedMTypes := []string{enum.Counter.String(), enum.Gauge.String()}
		if !slices.Contains(supportedMTypes, metricType) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var convMetricValue any
		var err error

		switch metricType {
		case enum.Counter.String():
			convMetricValue, err = strconv.ParseInt(metricValue, 0, 64)
		case enum.Gauge.String():
			convMetricValue, err = strconv.ParseFloat(metricValue, 64)
		}

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", enum.TextPlain.String())
		res.WriteHeader(http.StatusOK)
		h.memStorage.UpdateMetricValue(metricName, convMetricValue)
	}
}
