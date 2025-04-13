package config

import (
	"fmt"
	"net/url"

	"github.com/blackmou5e/grafana-annotator/pkg/errors"

	"github.com/spf13/viper"
)

type Config struct {
	GrafanaHost                string `mapstructure:"grafana_host"`
	GrafanaPort                string `mapstructure:"grafana_port"`
	GrafanaServiceAccountToken string `mapstructure:"sa_token"`
	Debug                      bool   `mapstructure:"debug"`
	Timeout                    int    `mapstructure:"timeout"`
}

func (c *Config) Validate() error {
	if c.GrafanaHost == "" {
		return errors.NewAppError(errors.ErrConfiguration, "Grafana host cannot be empty", nil)
	}

	if c.GrafanaServiceAccountToken == "" {
		return errors.NewAppError(errors.ErrConfiguration, "Service account token cannot be empty", nil)
	}

	_, err := url.Parse(fmt.Sprintf("http://%s:%s", c.GrafanaHost, c.GrafanaPort))
	if err != nil {
		return errors.NewAppError(errors.ErrConfiguration, "Invalid Grafana URL", err)
	}

	if c.Timeout == 0 {
		c.Timeout = 30 // Default timeout in seconds
	}

	return nil
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.gf-annotator")
	viper.AddConfigPath(".")

	viper.SetDefault("timeout", 30)
	viper.SetDefault("debug", false)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.NewAppError(errors.ErrConfiguration, "Failed to read config file", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.NewAppError(errors.ErrConfiguration, "Failed to parse config", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
