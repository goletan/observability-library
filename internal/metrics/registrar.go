package metrics

// Registrar defines an interface for entities that can be registered, providing a Register method that returns an error.
type Registrar interface {
	Register() error
}
