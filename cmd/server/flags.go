package main

import (
	"flag"
	"log"
	"net/url"
	"strings"

	"github.com/anatolev-max/metrics-tpl/config"
)

var options struct {
	flagRunAddr string
}

func parseFlags(c config.Config) {
	endpoint := c.Server.Schema + c.Server.Host + c.Server.Port
	flag.StringVar(&options.flagRunAddr, "a", endpoint, "address and port to run server")
	flag.Parse()

	if !strings.Contains(options.flagRunAddr, c.Server.Schema) {
		options.flagRunAddr = c.Server.Schema + options.flagRunAddr
	}

	validateFlags()
}

func validateFlags() {
	u, err := url.ParseRequestURI(options.flagRunAddr)
	if err != nil || u.Port() == "" {
		log.Fatal("Error while parsing command line arguments")
	}
}
