package types

import "time"

type ObservabilityConfig struct {
	Environment string `mapstructure:"environment"`

	Metrics struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"metrics"`

	Bulkhead struct {
		Capacity int           `mapstructure:"capacity"`
		Timeout  time.Duration `mapstructure:"timeout"`
	} `mapstructure:"bulkhead"`

	Logger struct {
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"logger"`

	Tracing struct {
		ServiceName   string  `mapstructure:"service_name"`
		Exporter      string  `mapstructure:"exporter"`
		SamplingRatio float64 `mapstructure:"sampling_ratio"`
	} `mapstructure:"tracing"`
}
