package config

import (
	config "github.com/goletan/config/pkg"
	"github.com/goletan/observability/internal/types"
)

func LoadObservabilityConfig() (*types.ObservabilityConfig, error) {
	cfg := &types.ObservabilityConfig{}
	if err := config.LoadConfig("Observability", cfg, nil); err != nil {
		return nil, err
	}

	return cfg, nil
}
