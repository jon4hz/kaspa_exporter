package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jon4hz/kaspa_exporter/metrics"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type KaspaCollector struct {
	Metrics   []metrics.Metrics
	Logger    log.Logger
	Namespace string
	client    *rpcclient.RPCClient
}

func New(logger log.Logger, namespace string, client *rpcclient.RPCClient) *KaspaCollector {
	c := &KaspaCollector{
		Namespace: namespace,
		Logger:    logger,
		client:    client,
	}
	c.Metrics = metrics.New(logger, namespace)
	return c
}

func (c *KaspaCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.Metrics {
		for _, d := range m.Desc() {
			ch <- d
		}
	}
}

func (c *KaspaCollector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.Metrics {
		if err := m.Collect(c.client); err != nil {
			level.Error(c.Logger).Log("msg", "Failed to collect metrics", "collector", m.String(), "err", err)
			continue
		}
		metrics, err := m.Get()
		if err != nil {
			level.Error(c.Logger).Log("msg", "Failed to get metrics", "collector", m.String(), "err", err)
			continue
		}
		for _, metric := range metrics {
			ch <- metric
		}
	}
}
