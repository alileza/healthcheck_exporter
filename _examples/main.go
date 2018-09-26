package main

import (
	"net/http"
	"time"

	healthcheck "github.com/alileza/healthcheck_exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	healthcheck.Register(healthcheck.Spec{
		Name:     "my-awesome-database",
		Interval: time.Second,
		Handle: func() error {
			return nil
		},
	})
	healthcheck.Run()

	http.ListenAndServe(":9000", promhttp.Handler())

}
