package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Listen        string `mapstructure:"listen"`
	TLSConfigFile string `mapstructure:"tls_config_file"`
	Namespace     string `mapstructure:"namespace"`
	KaspaURL      string `mapstructure:"kaspa_url"`
}

// Load loads the config file.
// It searches in the following locations:
//
// /etc/kaspa_exporter/config.yml,
// $HOME/.config/kaspa_exporter/config.yml,
// config.yml
//
// command arguments will overwrite the value from the config
func Load(path string) (cfg *Config, err error) {
	if path != "" {
		return load(path)
	}
	for _, f := range [4]string{
		".config.yml",
		"config.yml",
	} {
		cfg, err = load(f)
		if err != nil && os.IsNotExist(err) {
			err = nil
			continue
		} else if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
			err = nil
			continue
		}
	}
	if cfg == nil {
		return cfg, viper.Unmarshal(&cfg)
	}
	return
}

func load(file string) (cfg *Config, err error) {
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/kaspa_exporter")
	viper.AddConfigPath("/etc/kaspa_exporter/")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}
	viper.AutomaticEnv()
	return
}
