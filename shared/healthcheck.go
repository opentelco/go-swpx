package shared

// Healthcheck is the implementation for the protcolbuffer Healthcheck
type Healthcheck interface {
	Health() (bool, error)
}
