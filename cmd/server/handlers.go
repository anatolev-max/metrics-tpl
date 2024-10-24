package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var ms = NewMemStorage()
var msPointer = &ms

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var data []byte

	switch MetricType(req.PathValue("type")) {
	case Counter:
		data, _ = json.Marshal(msPointer.Counter)
	case Gauge:
		data, _ = json.Marshal(msPointer.Gauge)
	default:
		res.WriteHeader(http.StatusBadRequest)
		http.Error(res, "Metric type is not supported!", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(data); err != nil {
		panic(err)
	}
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if ContentType(req.Header.Get("Content-Type")) != TextPlain {
		http.Error(res, "Only text/plain are allowed!", http.StatusBadRequest)
		return
	}

	var metricValue interface{}
	strMetricValue := req.PathValue("value")

	switch MetricType(req.PathValue("type")) {
	case Counter:
		metricValue, _ = strconv.ParseInt(strMetricValue, 0, 64)
	case Gauge:
		metricValue, _ = strconv.ParseFloat(strMetricValue, 64)
	default:
		res.WriteHeader(http.StatusBadRequest)
		http.Error(res, "Metric type is not supported!", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	msPointer.Update(req.PathValue("name"), metricValue)
}
