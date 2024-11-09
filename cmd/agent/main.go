package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	. "github.com/anatolev-max/metrics-tpl/config"
	"github.com/anatolev-max/metrics-tpl/internal/enum"
	"github.com/anatolev-max/metrics-tpl/internal/storage"

	"github.com/go-resty/resty/v2"
)

func main() {
	c := NewConfig()
	parseFlags(c)

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := storage.NewMemStorage()
	fmt.Println("Running agent\nServer endpoint", options.flagRunAddr)
	diff := options.reportInterval % pollInterval

	var i uint = 0
	for {
		i++
		s.UpdateAgentData()

		if i%(options.reportInterval/pollInterval) == 0 {
			if diff != 0 {
				time.Sleep(time.Duration(diff) * time.Second)
			}
			updateServerData(s)
			i = 0
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func updateServerData(s storage.MemStorage) {
	sValue := reflect.ValueOf(s)

	for fieldIndex := 0; fieldIndex < sValue.NumField(); fieldIndex++ {
		metricType := sValue.Type().Field(fieldIndex).Name

		iter := sValue.FieldByName(metricType).MapRange()
		for iter.Next() {
			metricType = strings.ToLower(metricType)
			url := fmt.Sprintf(options.flagRunAddr+enum.UpdateEndpoint+"/%v/%v/%v", metricType, iter.Key(), iter.Value())
			sendRequest(url)
		}
	}
}

func sendRequest(url string) {
	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", enum.TextPlain).
		Post(url)

	if err != nil {
		log.Fatalln(err)
	}
}
