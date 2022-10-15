package mempool

import (
	"fmt"

	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	descs   []*prometheus.Desc
	entries float64
	orphans float64
}

func (c *Collector) String() string { return "mempool" }

func (c *Collector) Init(namespace string) error {
	c.descs = []*prometheus.Desc{
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "entries"), "Total amount of mempool entires", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "orphans"), "Amount of orphans in mempool", nil, nil),
	}
	return nil
}

func (c *Collector) Desc() []*prometheus.Desc {
	return c.descs
}

func (c *Collector) Collect(client *rpcclient.RPCClient) error {
	entries, err := client.GetMempoolEntries(true, false)
	if err != nil {
		return fmt.Errorf("failed to get mempool: %w", err)
	}
	c.entries = float64(len(entries.Entries))
	orphans := 0.0
	for _, entry := range entries.Entries {
		if entry.IsOrphan {
			orphans++
		}
	}
	c.orphans = orphans
	return nil
}

func (c *Collector) Get() ([]prometheus.Metric, error) {
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(c.Desc()[0], prometheus.GaugeValue, c.entries),
		prometheus.MustNewConstMetric(c.Desc()[1], prometheus.GaugeValue, c.orphans),
	}, nil
}
