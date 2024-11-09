package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/anatolev-max/metrics-tpl/config"
)

const pollInterval uint = 2
const reportInterval uint = 10

var options struct {
	flagRunAddr    string
	reportInterval uint
}

func parseFlags(c config.Config) {
	endpoint := c.Server.Host + c.Server.Port
	flag.StringVar(&options.flagRunAddr, "a", endpoint, "address and port for sending http requests")
	flag.UintVar(&options.reportInterval, "r", reportInterval, "frequency of sending metrics to the server")
	flag.Parse()

	validateFlags()
}

func validateFlags() {
	u, err := url.ParseRequestURI(options.flagRunAddr)
	if err != nil || u.Port() == "" {
		log.Fatal("Error while parsing command line arguments")
	}
}
