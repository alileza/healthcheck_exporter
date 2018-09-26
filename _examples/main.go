package main

import (
	"net/http"
	"time"

	healthcheck "github.com/alileza/healthcheck_exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	healthcheck.Register("my-awesome-database", time.Second, healthcheck.Handler(func() error {
		return nil
	}))
	healthcheck.Run()

	http.ListenAndServe(":9000", promhttp.Handler())
}
