package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/anatolev-max/metrics-tpl/cmd/common"
	"github.com/anatolev-max/metrics-tpl/cmd/storage"
	"github.com/go-chi/chi/v5"
)

func GetMainWebhook(ms storage.MemStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		data := struct {
			Title    string
			Counters map[string]int64
			Gauges   map[string]float64
		}{
			Title:    "Metrics-tpl",
			Counters: ms.Counter,
			Gauges:   ms.Gauge,
		}

		tmpl := template.Must(template.New("data").Parse(`<div>
				<h1>{{ .Title}}</h1>
			  	<ul>
					{{range $k, $v := .Counters}}
						<li>{{ $k }} - {{ $v }}</li>
					{{end}}
				</ul>
				<ul>
					{{range $k, $v := .Gauges}}
						<li>{{ $k }} - {{ $v }}</li>
					{{end}}
				</ul>
			</div>`))

		err := tmpl.Execute(res, data)
		if err != nil {
			panic(err)
		}
	}
}

func GetValueWebhook(ms storage.MemStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var data []byte
		metricType := strings.ToLower(chi.URLParam(req, "type"))
		metricName := strings.ToLower(chi.URLParam(req, "name"))

		supportedMTypes := []string{common.Counter, common.Gauge}
		if !slices.Contains(supportedMTypes, metricType) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		msVal := reflect.ValueOf(ms)

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

		res.Header().Set("Content-Type", common.TextPlain)
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
		urlPath := strings.TrimLeft(req.RequestURI, common.ServerHost+common.ServerPort)
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

		supportedMTypes := []string{common.Counter, common.Gauge}
		if req.Header.Get("Content-Type") != common.TextPlain || !slices.Contains(supportedMTypes, metricType) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var convMetricValue interface{}
		var err error

		switch metricType {
		case common.Counter:
			convMetricValue, err = strconv.ParseInt(metricValue, 0, 64)
		case common.Gauge:
			convMetricValue, err = strconv.ParseFloat(metricValue, 64)
		}

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", common.TextPlain)
		res.WriteHeader(http.StatusOK)
		ms.UpdateServerData(metricName, convMetricValue)
	}
}
