package types

import "time"

type ObservabilityConfig struct {
	Environment string `mapstructure:"environment"`

	Bulkhead struct {
		Capacity int           `mapstructure:"capacity"`
		Timeout  time.Duration `mapstructure:"timeout"`
	} `mapstructure:"bulkhead"`

	Logger struct {
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"logger-library"`

	Tracing struct {
		ServiceName   string  `mapstructure:"service_name"`
		Exporter      string  `mapstructure:"exporter"`
		SamplingRatio float64 `mapstructure:"sampling_ratio"`
	} `mapstructure:"tracing"`
}
