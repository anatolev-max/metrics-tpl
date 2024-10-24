package main

type MetricType string
type ContentType string

const (
	Counter = MetricType("counter")
	Gauge   = MetricType("gauge")
)

const (
	TextPlain = ContentType("text/plain")
)
