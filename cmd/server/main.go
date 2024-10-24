package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorageInterface interface {
	update(name string, value interface{})
}

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func (ms *MemStorage) update(name string, value interface{}) {
	switch vv := value.(type) {
	case int64:
		if _, exist := ms.Counter[name]; !exist {
			(*ms).Counter[name] = vv
		} else {
			(*ms).Counter[name] += vv
		}
	case float64:
		(*ms).Gauge[name] = vv
	}
}

var msPointer = &MemStorage{
	Counter: map[string]int64{},
	Gauge:   map[string]float64{},
}

const TextPlainContentType string = "text/plain"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{name}/{value}`, PostHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != TextPlainContentType {
		http.Error(res, "Only text/plain are allowed!", http.StatusBadRequest)
		return
	}

	queryParams := strings.Split(req.URL.RequestURI(), "/")
	if len(queryParams) != 5 {
		return
	}

	metricType := queryParams[1:][1]
	metricName := queryParams[1:][2]

	switch metricType {
	case "counter":
		res.WriteHeader(http.StatusOK)
		metricValue, _ := strconv.ParseInt(queryParams[1:][3], 0, 64)
		msPointer.update(metricName, metricValue)
	case "gauge":
		res.WriteHeader(http.StatusOK)
		metricValue, _ := strconv.ParseFloat(queryParams[1:][3], 64)
		msPointer.update(metricName, metricValue)
	default:
		res.WriteHeader(http.StatusBadRequest)
		http.Error(res, "Metric type is not supported!", http.StatusBadRequest)
	}
}
