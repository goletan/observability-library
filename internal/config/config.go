package config

import (
	config "github.com/goletan/config/pkg"
	"github.com/goletan/observability/internal/types"
	"go.uber.org/zap"
)

var cfg types.ObservabilityConfig

func LoadObservabilityConfig(logger *zap.Logger) (*types.ObservabilityConfig, error) {
	if err := config.LoadConfig("Observability", &cfg, logger); err != nil {
		logger.Error(
			"Failed to load events configuration",
			zap.Error(err),
			zap.Any("context", map[string]interface{}{"step": "config loading"}),
		)
		return nil, err
	}

	return &cfg, nil
}
