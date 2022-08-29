package info

import (
	"fmt"

	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	descs       []*prometheus.Desc
	synced      float64
	mempoolSize float64
	difficulty  float64
	hashrate    float64
}

func (c *Collector) String() string { return "info" }

func (c *Collector) Init(namespace string) error {
	c.descs = []*prometheus.Desc{
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "sync"), "Get the sync status", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "mempool_size"), "Get the mempool size", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "difficulty"), "Get the network difficulty", nil, nil),
		prometheus.NewDesc(prometheus.BuildFQName(namespace, c.String(), "hashrate"), "Get the network hashrate", nil, nil),
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
	if info.IsSynced {
		c.synced = 1
	}
	c.mempoolSize = float64(info.MempoolSize)

	blockDAG, err := client.GetBlockDAGInfo()
	if err != nil {
		return fmt.Errorf("failed to get block dag info: %w", err)
	}
	c.difficulty = blockDAG.Difficulty

	hr, err := client.EstimateNetworkHashesPerSecond(blockDAG.TipHashes[0], 1000)
	if err != nil {
		return fmt.Errorf("failed to get network hashrate: %w", err)
	}
	c.hashrate = float64(hr.NetworkHashesPerSecond)
	return nil
}

func (c *Collector) Get() ([]prometheus.Metric, error) {
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(c.Desc()[0], prometheus.GaugeValue, c.synced),
		prometheus.MustNewConstMetric(c.Desc()[1], prometheus.GaugeValue, c.mempoolSize),
		prometheus.MustNewConstMetric(c.Desc()[2], prometheus.GaugeValue, c.difficulty),
		prometheus.MustNewConstMetric(c.Desc()[3], prometheus.GaugeValue, c.hashrate),
	}, nil
}
