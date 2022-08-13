package peers

import (
	"fmt"

	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	descs    []*prometheus.Desc
	total    float64
	inbound  float64
	outbound float64
}

func (c *Collector) String() string { return "peers" }

func (c *Collector) Init(namespace string) error {
	c.descs = []*prometheus.Desc{
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "count"), "Total of connected peers", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "inbound_count"), "Total of inbound peers", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "outbound_count"), "Total of outbound peers", nil, nil),
	}
	return nil
}

func (c *Collector) Desc() []*prometheus.Desc {
	return c.descs
}

func (c *Collector) Collect(client *rpcclient.RPCClient) error {
	peers, err := client.GetConnectedPeerInfo()
	if err != nil {
		return fmt.Errorf("failed to get peers: %w", err)
	}
	c.total = float64(len(peers.Infos))
	for _, i := range peers.Infos {
		if i.IsIBDPeer {
			c.inbound++
		}
		if i.IsOutbound {
			c.outbound++
		}
	}
	return nil
}

func (c *Collector) Get() ([]prometheus.Metric, error) {
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(c.Desc()[0], prometheus.GaugeValue, c.total),
		prometheus.MustNewConstMetric(c.Desc()[1], prometheus.GaugeValue, c.inbound),
		prometheus.MustNewConstMetric(c.Desc()[2], prometheus.GaugeValue, c.outbound),
	}, nil
}
