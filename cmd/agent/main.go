package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/anatolev-max/metrics-tpl/cmd/common"
	"github.com/anatolev-max/metrics-tpl/cmd/storage"
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
	memStorage := storage.NewMemStorage()

	for {
		memStorage.UpdateAgentData()
		if memStorage.Counter[common.PollCounter]%(reportInterval/pollInterval) == 0 {
			updateServerData(memStorage)
		}

		time.Sleep(pollInterval * time.Second)
	}
}

func updateServerData(ms storage.MemStorage) {
	msVal := reflect.ValueOf(ms)

	for fieldIndex := 0; fieldIndex < msVal.NumField(); fieldIndex++ {
		metricType := msVal.Type().Field(fieldIndex).Name

		iter := msVal.FieldByName(metricType).MapRange()
		for iter.Next() {
			metricType = strings.ToLower(metricType)
			url := fmt.Sprintf(common.UpdateFullEndpoint+"%v/%v/%v", metricType, iter.Key(), iter.Value())
			sendRequest(url)
		}
	}
}

func sendRequest(url string) {
	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", common.TextPlain).
		Post(url)

	if err != nil {
		log.Fatalln(err)
	}
}
