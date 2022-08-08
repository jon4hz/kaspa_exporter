package metrics

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jon4hz/kaspa_exporter/metrics/info"
	"github.com/jon4hz/kaspa_exporter/metrics/mempool"
	"github.com/jon4hz/kaspa_exporter/metrics/peers"
	"github.com/jon4hz/kaspa_exporter/metrics/supply"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

var availableMetrics = []Metrics{
	&info.Collector{},
	&mempool.Collector{},
	&peers.Collector{},
	&supply.Collector{},
}

type Metrics interface {
	Init(namespace string) error
	Collect(*rpcclient.RPCClient) error
	Get() ([]prometheus.Metric, error)
	Desc() []*prometheus.Desc
	String() string
}

func New(logger log.Logger, namespace string) []Metrics {
	metrics := make([]Metrics, 0)
	for _, m := range availableMetrics {
		if err := m.Init(namespace); err != nil {
			level.Error(logger).Log("msg", "Failed to initialize metrics", "collector", m.String(), "err", err)
			continue
		}
		metrics = append(metrics, m)
	}
	return metrics
}
