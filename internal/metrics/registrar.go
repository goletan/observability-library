package metrics

type MetricsRegistrar interface {
	Register() error
}
