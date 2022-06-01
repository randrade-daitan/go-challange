package api

// An API server
type Api interface {

	// Begin advertising the service at a given port.
	StartServing(port int) error
}
