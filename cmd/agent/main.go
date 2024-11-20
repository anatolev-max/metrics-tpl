package main

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/anatolev-max/metrics-tpl/config"
	"github.com/anatolev-max/metrics-tpl/internal/enum"
	"github.com/anatolev-max/metrics-tpl/internal/storage"

	"github.com/go-resty/resty/v2"
)

func main() {
	c := config.NewConfig()
	parseFlags(c)

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := storage.NewMemStorage()
	fmt.Println("Running agent\nServer endpoint", options.runAddr)
	intervals := []uint{options.pollInterval, options.reportInterval}
	maxInterval := slices.Max(intervals)
	minInterval := slices.Min(intervals)
	diff := maxInterval % minInterval
	reportFirst := options.reportInterval > options.pollInterval

	var i uint
	for {
		i++
		sleep(minInterval)
		reportOrPoll(!reportFirst, s)

		if i%(maxInterval/minInterval) == 0 {
			if diff != 0 {
				sleep(diff)
			}
			reportOrPoll(reportFirst, s)
		}
	}
}

func reportOrPoll(report bool, s *storage.MemStorage) {
	if report {
		sendDataToServer(s)
	} else {
		s.UpdateAgentData()
	}
}

func sendDataToServer(s *storage.MemStorage) {
	sValue := reflect.ValueOf(s).Elem()

	for fieldIndex := 0; fieldIndex < sValue.NumField(); fieldIndex++ {
		metricType := sValue.Type().Field(fieldIndex).Name

		iter := sValue.FieldByName(metricType).MapRange()
		for iter.Next() {
			metricType = strings.ToLower(metricType)
			url := fmt.Sprintf(options.runAddr+enum.UpdateEndpoint.String()+"/%v/%v/%v", metricType, iter.Key(), iter.Value())
			sendRequest(url)
		}
	}
}

func sendRequest(url string) {
	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", enum.TextPlain.String()).
		Post(url)

	if err != nil {
		log.Fatalln(err)
	}
}

func sleep(seconds uint) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
