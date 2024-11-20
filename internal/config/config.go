package config

import (
	config "github.com/goletan/config/pkg"
	"github.com/goletan/observability/internal/types"
)

var cfg types.ObservabilityConfig

func LoadObservabilityConfig() (*types.ObservabilityConfig, error) {
	if err := config.LoadConfig("Observability", &cfg, nil); err != nil {
		return nil, err
	}

	return &cfg, nil
}
