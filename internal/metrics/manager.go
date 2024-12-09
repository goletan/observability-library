package metrics

import (
	"github.com/goletan/observability/shared/errors"
	"github.com/goletan/observability/shared/logger"
)

// Manager manages the registration and initialization of metrics components.
type Manager struct {
	registrars []Registrar
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Register(registrar Registrar) {
	m.registrars = append(m.registrars, registrar)
}

func (m *Manager) Init(log *logger.ZapLogger) error {
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
