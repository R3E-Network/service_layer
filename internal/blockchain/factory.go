package blockchain

import (
	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/pkg/logger"
)

// NeoBlockchainClientFactory creates Neo N3 blockchain clients
type NeoBlockchainClientFactory struct {
	config *config.NeoConfig
	logger *logger.Logger
}

// NewNeoBlockchainClientFactory creates a new Neo N3 blockchain client factory
func NewNeoBlockchainClientFactory(config *config.NeoConfig, logger *logger.Logger) *NeoBlockchainClientFactory {
	return &NeoBlockchainClientFactory{
		config: config,
		logger: logger,
	}
}

// NewClient implements the BlockchainClientFactory interface
func (f *NeoBlockchainClientFactory) NewClient() (BlockchainClient, error) {
	// Create the underlying client
	client, err := NewClient(f.config, f.logger)
	if err != nil {
		return nil, err
	}

	// Return a client adapter that implements the BlockchainClient interface
	return NewClientAdapter(client), nil
}

// NewMockClient implements the BlockchainClientFactory interface
func (f *NeoBlockchainClientFactory) NewMockClient() BlockchainClient {
	// Create a mock client for testing
	return &MockBlockchainClient{}
}

// Create is a convenience method that calls NewClient
func (f *NeoBlockchainClientFactory) Create() (BlockchainClient, error) {
	return f.NewClient()
}

// Initialize the factory
var factory *NeoBlockchainClientFactory

func init() {
	// Initialize the factory with default config and logger
	// In a real application, this would be configured via dependency injection
	factory = NewNeoBlockchainClientFactory(nil, nil)
}
