package config

import (
	config "github.com/goletan/config/pkg"
	"github.com/goletan/observability/internal/logger"
	"github.com/goletan/observability/internal/types"
	"go.uber.org/zap"
)

// LoadObservabilityConfig loads the observability configuration settings from a predefined source.
// It returns the loaded ObservabilityConfig struct and an error if the configuration loading fails.
func LoadObservabilityConfig(log *logger.ZapLogger) (*types.ObservabilityConfig, error) {
	var cfg types.ObservabilityConfig
	if err := config.LoadConfig("Observability", &cfg, nil); err != nil {
		log.Error("Failed to load observability configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
