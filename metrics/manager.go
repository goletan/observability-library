// /observability/metrics/manager.go
package metrics

import (
	"log"
)

type MetricsManager struct {
	registrars []MetricsRegistrar
}

var manager *MetricsManager

func NewManager() *MetricsManager {
	if manager == nil {
		manager = &MetricsManager{}
	}

	return manager
}

func (m *MetricsManager) Register(registrar MetricsRegistrar) {
	m.registrars = append(m.registrars, registrar)
}

func (m *MetricsManager) Init() {
	for _, registrar := range m.registrars {
		if err := registrar.Register(); err != nil {
			log.Printf("Warning: failed to register metrics: %v", err)
		}
	}
}
