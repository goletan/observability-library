// /observability/metrics/registrar.go
package metrics

type MetricsRegistrar interface {
	Register() error
}
