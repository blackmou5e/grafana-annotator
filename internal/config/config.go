package config

import (
	"net/url"

	"github.com/blackmou5e/grafana-annotator/pkg/errors"

	"github.com/spf13/viper"
)

type Config struct {
	GrafanaURL                 string `mapstructure:"annotator_grafana_url"`
	GrafanaServiceAccountToken string `mapstructure:"annotator_sa_token"`
	Debug                      bool   `mapstructure:"annotator_debug"`
	Timeout                    int    `mapstructure:"annotator_timeout"`
}

func (c *Config) Validate() error {
	if c.GrafanaURL == "" {
		return errors.NewAppError(errors.ErrConfiguration, "Grafana url cannot be empty", nil)
	}

	if c.GrafanaServiceAccountToken == "" {
		return errors.NewAppError(errors.ErrConfiguration, "Service account token cannot be empty", nil)
	}

	_, err := url.Parse(c.GrafanaURL)
	if err != nil {
		return errors.NewAppError(errors.ErrConfiguration, "Invalid Grafana URL", err)
	}

	if c.Timeout == 0 {
		c.Timeout = 30 // Default timeout in seconds
	}

	return nil
}

func SetDefaultConfig() (*Config, error) {
	viper.SetDefault("annotator_timeout", 30)
	viper.SetDefault("annotator_debug", false)
	viper.SetDefault("annotator_grafana_url", "http://localhost:3000")

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.NewAppError(errors.ErrConfiguration, "Failed to create default config", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadConfigFromFile(cfg *Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.grafana-annotator")
	viper.AddConfigPath("$XDG_CONFIG_HOME/grafana-annotator/config")
	viper.AddConfigPath("grafana-annotator")

	if err := viper.ReadInConfig(); err != nil {
		return errors.NewAppError(errors.ErrConfiguration, "Failed to read config file", err)
	}

	if err := unmarshalViperConfig(cfg); err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	return nil
}

func LoadConfigFromEnv(cfg *Config) error {
	viper.SetEnvPrefix("annotator")

	viper.AutomaticEnv()

	if err := unmarshalViperConfig(cfg); err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	return nil
}

func unmarshalViperConfig(cfg *Config) error {
	if err := viper.Unmarshal(cfg); err != nil {
		return errors.NewAppError(errors.ErrConfiguration, "Failed to save viper config to struct", err)
	}

	return nil
}
