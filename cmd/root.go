package cmd

import (
	"context"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jon4hz/kaspa_exporter/collector"
	"github.com/jon4hz/kaspa_exporter/internal/config"
	"github.com/jon4hz/kaspa_exporter/internal/version"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootFlags struct {
	config string
}

var rootCmd = &cobra.Command{
	Use:   "kaspa_exporter",
	Short: "Prometheus exporter for Kaspa",
	Run:   root,
}

func init() {
	rootCmd.AddCommand(versionCmd, manCmd)

	rootCmd.Flags().StringVarP(&rootFlags.config, "config", "c", "", "path to the config file")

	rootCmd.Flags().StringP("listen", "l", "0.0.0.0:5000", "Metrics endpoint")
	rootCmd.Flags().StringP("tls-config-file", "t", "", "TLS configuration file")
	rootCmd.Flags().StringP("namespace", "n", "kaspa", "Namespace")
	rootCmd.Flags().StringP("kaspa-url", "u", "", "Kaspa URL")

	viper.BindPFlag("listen", rootCmd.Flags().Lookup("listen"))                   // nolint: errcheck
	viper.BindPFlag("tls_config_file", rootCmd.Flags().Lookup("tls-config-file")) // nolint: errcheck
	viper.BindPFlag("namespace", rootCmd.Flags().Lookup("namespace"))             // nolint: errcheck
	viper.BindPFlag("kaspa_url", rootCmd.Flags().Lookup("kaspa-url"))             // nolint: errcheck
}

func root(cmd *cobra.Command, args []string) {
	cfg, err := config.Load(rootFlags.config)
	if err != nil {
		golog.Fatalln(err)
	}

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting kaspa_exporter", "version", version.Version) // nolint: errcheck

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		probeHandler(w, req, cfg, logger)
	})

	server := &http.Server{Addr: cfg.Listen, ReadHeaderTimeout: 5 * time.Second}
	if err := web.ListenAndServe(
		server,
		&web.FlagConfig{
			WebListenAddresses: getPointer([]string{cfg.Listen}),
			WebSystemdSocket:   getPointer(false),
			WebConfigFile:      &cfg.TLSConfigFile,
		},
		logger,
	); err != nil {
		level.Error(logger).Log("msg", "Failed to start the server", "err", err) // nolint: errcheck
		os.Exit(1)
	}
}

func getPointer[T any](v T) *T {
	return &v
}

func Execute() error {
	return rootCmd.Execute()
}

func probeHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config, logger log.Logger) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	r = r.WithContext(ctx)

	client, err := rpcclient.NewRPCClient(cfg.KaspaURL)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to start the server", "err", err) // nolint: errcheck
		return
	}
	defer client.Close()

	registry := prometheus.NewPedanticRegistry()

	col := collector.New(logger, cfg.Namespace, client)

	registry.MustRegister(col)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

var manCmd = &cobra.Command{
	Use:                   "man",
	Short:                 "generates the manpages",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		manPage, err := mcobra.NewManPage(1, rootCmd)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
		return err
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Date: %s\n", version.Date)
		fmt.Printf("Build by: %s\n", version.BuiltBy)
	},
}
