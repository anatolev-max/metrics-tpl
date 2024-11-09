package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/anatolev-max/metrics-tpl/config"
)

var options struct {
	flagRunAddr string
}

func parseFlags(c config.Config) {
	endpoint := c.Server.Host + c.Server.Port
	flag.StringVar(&options.flagRunAddr, "a", endpoint, "address and port to run server")
	flag.Parse()

	validateFlags()
}

func validateFlags() {
	u, err := url.ParseRequestURI(options.flagRunAddr)
	if err != nil || u.Port() == "" {
		log.Fatal("Error while parsing command line arguments")
	}
}
