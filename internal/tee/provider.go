package tee

import (
	"context"

	"github.com/R3E-Network/service_layer/internal/models"
)

// Provider defines the interface for TEE providers
type Provider interface {
	// Initialize initializes the provider
	Initialize() error

	// ExecuteFunction executes a function in the TEE
	ExecuteFunction(ctx context.Context, function *models.Function, params map[string]interface{}, secrets map[string]string) (*models.ExecutionResult, error)

	// StoreSecret securely stores a secret in the TEE
	StoreSecret(ctx context.Context, secret *models.Secret) error

	// GetSecret retrieves a secret from the TEE
	GetSecret(ctx context.Context, userID int, secretName string) (string, error)

	// DeleteSecret deletes a secret from the TEE
	DeleteSecret(ctx context.Context, userID int, secretName string) error

	// GetAttestation gets an attestation report from the TEE
	GetAttestation(ctx context.Context) ([]byte, error)

	// Close closes the provider
	Close() error
}
