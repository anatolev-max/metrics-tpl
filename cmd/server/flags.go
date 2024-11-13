package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/anatolev-max/metrics-tpl/config"
)

var options struct {
	flagRunAddr string
}

func parseFlags(c config.Config) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		options.flagRunAddr = envRunAddr
	} else {
		hp := c.Server.Host + c.Server.Port
		flag.StringVar(&options.flagRunAddr, "a", hp, "address and port to run server")
	}

	flag.Parse()
	options.flagRunAddr = c.Server.Scheme + options.flagRunAddr

	validateFlags()
}

func validateFlags() {
	u, err := url.ParseRequestURI(options.flagRunAddr)
	if err != nil || u.Port() == "" {
		log.Fatal("Error while parsing command line arguments")
	}
}
