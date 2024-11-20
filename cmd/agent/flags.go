package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/anatolev-max/metrics-tpl/config"
)

const pollInterval uint = 2
const reportInterval uint = 10

var options struct {
	runAddr        string
	pollInterval   uint
	reportInterval uint
}

func parseFlags(c config.Config) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		options.runAddr = envRunAddr
	} else {
		hp := c.Server.Host + c.Server.Port
		flag.StringVar(&options.runAddr, "a", hp, "address and port for sending http requests")
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		v, _ := strconv.ParseInt(envReportInterval, 10, 32)
		options.reportInterval = uint(v)
	} else {
		flag.UintVar(&options.reportInterval, "r", reportInterval, "frequency of sending metrics to the server")
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		v, _ := strconv.ParseInt(envPollInterval, 10, 32)
		options.pollInterval = uint(v)
	} else {
		flag.UintVar(&options.pollInterval, "p", pollInterval, "frequency of polling metrics from the runtime package")
	}

	flag.Parse()
	options.runAddr = c.Server.Scheme + options.runAddr

	validateFlags()
}

func validateFlags() {
	u, err := url.ParseRequestURI(options.runAddr)
	if err != nil || u.Port() == "" || options.pollInterval < 1 || options.reportInterval < 1 {
		log.Fatal("Error while parsing command line arguments")
	}
}
