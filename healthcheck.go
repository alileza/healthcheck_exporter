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
	specs   []Spec
	tickers []*time.Ticker
}

func New() *Checker {
	return &Checker{}
}

type Spec struct {
	Name     string
	Interval time.Duration
	Handle   Handler
}

type Handler func() error

func Register(s Spec) {
	baseChecker.Register(s)
}

func (c *Checker) Register(s Spec) {
	c.specs = append(c.specs, s)
}

func Run() {
	baseChecker.Run()
}

func (c *Checker) Run() {
	for _, spec := range c.specs {
		t := func(spec Spec) *time.Ticker {
			ticker := time.NewTicker(spec.Interval)
			go func() {
				for range ticker.C {
					t := time.Now()

					err := spec.Handle()

					healthLatency.WithLabelValues(spec.Name).Set(time.Since(t).Seconds())
					if err != nil {
						healthStatus.WithLabelValues(spec.Name).Set(0)
						continue
					}
					healthStatus.WithLabelValues(spec.Name).Set(1)
				}
			}()
			return ticker
		}(spec)

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
