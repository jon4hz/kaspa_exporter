package main

import (
	"os"

	"github.com/jon4hz/kaspa_exporter/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
