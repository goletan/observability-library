// /observability/internal/metrics/manager.go
package metrics

import (
	"fmt"
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

func (m *MetricsManager) Init() error {
	for _, registrar := range m.registrars {
		if err := registrar.Register(); err != nil {
			return fmt.Errorf("failed to register metrics: %w", err)
		}
	}
	return nil
}
