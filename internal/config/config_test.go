package config_test

import (
	"testing"

	_ "github.com/jon4hz/kaspa_exporter/cmd"
	"github.com/jon4hz/kaspa_exporter/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := config.Load("testdata/config.yml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "127.0.0.1:5000", cfg.Listen)
	assert.Equal(t, "kaspa.io:16110", cfg.KaspaURL)
	assert.Equal(t, "", cfg.TLSConfigFile)
	assert.Equal(t, "kaspa", cfg.Namespace)
}

func TestLoadNoFileConfig(t *testing.T) {
	cfg, err := config.Load("")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

}

func TestLoadInvalidFile(t *testing.T) {
	_, err := config.Load("testdata/invalid.txt")
	assert.Error(t, err)
}
