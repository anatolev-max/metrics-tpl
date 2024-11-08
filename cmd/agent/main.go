package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/anatolev-max/metrics-tpl/internal/config"
	"github.com/anatolev-max/metrics-tpl/internal/storage"
	"github.com/go-resty/resty/v2"
)

const pollInterval = 2
const reportInterval = 10

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := storage.NewMemStorage()

	for {
		s.UpdateAgentData()
		if s.Counter[config.PollCounter]%(reportInterval/pollInterval) == 0 {
			updateServerData(s)
		}

		time.Sleep(pollInterval * time.Second)
	}
}

func updateServerData(s storage.MemStorage) {
	sValue := reflect.ValueOf(s)

	for fieldIndex := 0; fieldIndex < sValue.NumField(); fieldIndex++ {
		metricType := sValue.Type().Field(fieldIndex).Name

		iter := sValue.FieldByName(metricType).MapRange()
		for iter.Next() {
			metricType = strings.ToLower(metricType)
			url := fmt.Sprintf(config.UpdateFullEndpoint+"%v/%v/%v", metricType, iter.Key(), iter.Value())
			sendRequest(url)
		}
	}
}

func sendRequest(url string) {
	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", config.TextPlain).
		Post(url)

	if err != nil {
		log.Fatalln(err)
	}
}
