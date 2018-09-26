package healthcheck

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "healthcheck"

var (
	baseChecker = New()

	healthStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "up",
		Help:      "A gauge of healthcheck up status",
	}, []string{"name"})

	healthLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "latency",
		Help:      "A gauge of healthcheck latency",
	}, []string{"name"})
)

func init() {
	prometheus.MustRegister(healthStatus, healthLatency)
}

type Checker struct {
	specs   []spec
	tickers []*time.Ticker
}

func New() *Checker {
	return &Checker{}
}

type spec struct {
	Name     string
	Interval time.Duration
	Handle   Handler
}

type Handler func() error

func Register(name string, interval time.Duration, handle Handler) {
	baseChecker.Register(name, interval, handle)
}

func (c *Checker) Register(name string, interval time.Duration, handle Handler) {
	c.specs = append(c.specs, spec{
		Name:     name,
		Interval: interval,
		Handle:   handle,
	})
}

func Run() {
	baseChecker.Run()
}

func (c *Checker) Run() {
	for _, s := range c.specs {
		t := func(s spec) *time.Ticker {
			ticker := time.NewTicker(s.Interval)
			go func() {
				for range ticker.C {
					t := time.Now()

					err := s.Handle()

					healthLatency.WithLabelValues(s.Name).Set(time.Since(t).Seconds())
					if err != nil {
						healthStatus.WithLabelValues(s.Name).Set(0)
						continue
					}
					healthStatus.WithLabelValues(s.Name).Set(1)
				}
			}()
			return ticker
		}(s)

		c.tickers = append(c.tickers, t)
	}
}

func Close() {
	baseChecker.Close()
}

func (c *Checker) Close() {
	for _, t := range c.tickers {
		t.Stop()
	}
}
