package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jon4hz/kaspa_exporter/collector"
	"github.com/jon4hz/kaspa_exporter/internal/version"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
)

const (
	listenAddress = "0.0.0.0:5000"
	tlsConfigFile = ""
	namespace     = "kaspa"
)

var kaspaURL = os.Args[1]

func main() {

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting kaspa_exporter", "version", version.Version)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		probeHandler(w, req, logger, nil)
	})

	server := &http.Server{Addr: listenAddress}
	if err := web.ListenAndServe(server, tlsConfigFile, logger); err != nil {
		level.Error(logger).Log("msg", "Failed to start the server", "err", err)
		os.Exit(1)
	}

}

func probeHandler(w http.ResponseWriter, r *http.Request, logger log.Logger, _ *rpcclient.RPCClient) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	r = r.WithContext(ctx)

	registry := prometheus.NewPedanticRegistry()

	client, err := rpcclient.NewRPCClient(kaspaURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	col := collector.New(logger, namespace, client)

	registry.MustRegister(col)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
