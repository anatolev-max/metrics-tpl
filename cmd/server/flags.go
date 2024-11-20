package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/anatolev-max/metrics-tpl/config"
)

var options struct {
	runAddr string
}

func parseFlags(c config.Config) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		options.runAddr = envRunAddr
	} else {
		hp := c.Server.Host + c.Server.Port
		flag.StringVar(&options.runAddr, "a", hp, "address and port to run server")
	}

	flag.Parse()
	options.runAddr = c.Server.Scheme + options.runAddr

	validateFlags()
}

func validateFlags() {
	u, err := url.ParseRequestURI(options.runAddr)
	if err != nil || u.Port() == "" {
		log.Fatal("Error while parsing command line arguments")
	}
}
