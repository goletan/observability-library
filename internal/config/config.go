package config

import (
	"github.com/goletan/config-library/pkg"
	logger "github.com/goletan/logger-library/pkg"
	"github.com/goletan/observability/internal/types"
	"go.uber.org/zap"
)

var cfg types.ObservabilityConfig

// LoadObservabilityConfig loads the observability-library configuration settings from a predefined source.
// It returns the loaded ObservabilityConfig struct and an error if the configuration loading fails.
func LoadObservabilityConfig(log *logger.ZapLogger) (*types.ObservabilityConfig, error) {
	if err := config.LoadConfig("Observability", &cfg, log); err != nil {
		log.Error("Failed to load observability-library configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
