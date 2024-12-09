package metrics

import (
	"github.com/goletan/observability/shared/errors"
	"github.com/goletan/observability/shared/logger"
)

type MetricsManager struct {
	registrars []MetricsRegistrar
}

func NewManager() *MetricsManager {
	return &MetricsManager{}
}

func (m *MetricsManager) Register(registrar MetricsRegistrar) {
	m.registrars = append(m.registrars, registrar)
}

func (m *MetricsManager) Init(log *logger.ZapLogger) error {
	for _, registrar := range m.registrars {
		if err := registrar.Register(); err != nil {
			return errors.WrapError(log, err, "failed to register metrics", 2001, map[string]interface{}{
				"component": "metrics",
				"endpoint":  "register",
			})
		}
	}
	return nil
}
