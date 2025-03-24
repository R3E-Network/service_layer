package common

import "context"

// Service defines a common interface for all services in the system
type Service interface {
	// Start starts the service
	Start(ctx context.Context) error

	// Stop stops the service
	Stop() error

	// Name returns the service name
	Name() string
}
