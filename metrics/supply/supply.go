package supply

import (
	"fmt"

	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	descs  []*prometheus.Desc
	supply float64
}

func (c *Collector) String() string { return "supply" }

func (c *Collector) Init(namespace string) error {
	c.descs = []*prometheus.Desc{
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "circulating"), "Circulating supply", nil, nil),
	}
	return nil
}

func (c *Collector) Desc() []*prometheus.Desc {
	return c.descs
}

func (c *Collector) Collect(client *rpcclient.RPCClient) error {
	supply, err := client.GetCoinSupply()
	if err != nil {
		return fmt.Errorf("failed to get supply: %w", err)
	}
	c.supply = float64(supply.CirculatingSompi)
	return nil
}

func (c *Collector) Get() ([]prometheus.Metric, error) {
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(c.Desc()[0], prometheus.GaugeValue, c.supply),
	}, nil
}
