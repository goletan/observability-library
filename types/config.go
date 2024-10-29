// /observability/types/config.go
package types

import "time"

// Holds the configuration for the logger.
type ObservabilityConfig struct {
	Bulkhead struct {
		Capacity int           `mapstructure:"capacity"`
		Timeout  time.Duration `mapstructure:"timeout"`
	} `mapstructure:"bulkhead"`

	Logger struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`

	Tracing struct {
		ServiceName   string  `mapstructure:"service_name"`
		Exporter      string  `mapstructure:"exporter"`
		SamplingRatio float64 `mapstructure:"sampling_ratio"`
	} `mapstructure:"tracing"`
}
