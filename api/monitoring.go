package api

import (
	"regexp"

	"github.com/1755/bookstore-api/api/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type MonitoringConfig struct {
	Address string `validate:"required,ip"`
	Port    int    `validate:"required,min=1,max=65535"`
}

func NewMonitoring() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(
				collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/sched/latencies:seconds")},
			),
		),
		middlewares.HttpRequestsTotal,
		middlewares.HttpRequestDuration,
	)

	return registry
}
