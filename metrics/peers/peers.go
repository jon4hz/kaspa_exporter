package peers

import (
	"fmt"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	descs []*prometheus.Desc
	peers *appmessage.GetConnectedPeerInfoResponseMessage
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
	c.peers = peers
	return nil
}

func (c *Collector) Get() ([]prometheus.Metric, error) {
	var incoming, outgoing int
	for _, i := range c.peers.Infos {
		if i.IsIBDPeer {
			incoming++
		}
		if i.IsOutbound {
			outgoing++
		}
	}

	return []prometheus.Metric{
		prometheus.MustNewConstMetric(c.Desc()[0], prometheus.GaugeValue, float64(len(c.peers.Infos))),
		prometheus.MustNewConstMetric(c.Desc()[1], prometheus.GaugeValue, float64(incoming)),
		prometheus.MustNewConstMetric(c.Desc()[2], prometheus.GaugeValue, float64(outgoing)),
	}, nil
}
