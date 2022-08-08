package info

import (
	"fmt"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	descs []*prometheus.Desc
	info  *appmessage.GetInfoResponseMessage
}

func (c *Collector) String() string { return "info" }

func (c *Collector) Init(namespace string) error {
	c.descs = []*prometheus.Desc{
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "sync"), "Get the sync status", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "mempool_size"), "Get the mempool size", nil, nil),
	}
	return nil
}

func (c *Collector) Desc() []*prometheus.Desc {
	return c.descs
}

func (c *Collector) Collect(client *rpcclient.RPCClient) error {
	info, err := client.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get info: %w", err)
	}
	c.info = info
	return nil
}

func (c *Collector) Get() ([]prometheus.Metric, error) {
	var syncVal float64
	if c.info.IsSynced {
		syncVal = 1
	}
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(c.Desc()[0], prometheus.GaugeValue, syncVal),
		prometheus.MustNewConstMetric(c.Desc()[1], prometheus.GaugeValue, float64(c.info.MempoolSize)),
	}, nil
}
