package blockchain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/pkg/logger"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

// ClientFactory is an implementation of BlockchainClientFactory
type ClientFactory struct {
	config *config.Config
	logger *logger.Logger
}

// NewClientFactory creates a new blockchain client factory
func NewClientFactory(cfg *config.Config, logger *logger.Logger) *ClientFactory {
	return &ClientFactory{
		config: cfg,
		logger: logger,
	}
}

// NewClient creates a new blockchain client
func (f *ClientFactory) NewClient() (BlockchainClient, error) {
	// Configure nodes from config
	var nodes []NodeConfig
	for _, nodeURL := range f.config.Neo.Nodes {
		nodes = append(nodes, NodeConfig{
			URL:    nodeURL.StringURL(),
			Weight: 1.0, // Default equal weight for all nodes
		})
	}

	// Create a real blockchain client
	client, err := NewClient(&f.config.Neo, f.logger, nodes)
	if err != nil {
		return nil, err
	}

	// Create a blockchain client adapter that implements the BlockchainClient interface
	return &blockchainClientAdapter{
		client: client,
	}, nil
}

// blockchainClientAdapter implements the BlockchainClient interface by delegating to the Neo client
type blockchainClientAdapter struct {
	client *Client
}

// GetBlockHeight implements BlockchainClient interface
func (b *blockchainClientAdapter) GetBlockHeight() (uint32, error) {
	return b.client.GetBlockHeight()
}

// GetBlock implements BlockchainClient interface
func (b *blockchainClientAdapter) GetBlock(height uint32) (interface{}, error) {
	block, err := b.client.GetBlock(height)
	return block, err
}

// GetTransaction implements BlockchainClient interface
func (b *blockchainClientAdapter) GetTransaction(hash string) (interface{}, error) {
	tx, err := b.client.GetTransaction(hash)
	return tx, err
}

// SendTransaction implements BlockchainClient interface
func (b *blockchainClientAdapter) SendTransaction(tx interface{}) (string, error) {
	// Try to convert to the expected type
	txObj, ok := tx.(*transaction.Transaction)
	if !ok {
		return "", fmt.Errorf("invalid transaction type: expected *transaction.Transaction, got %T", tx)
	}
	return b.client.SendTransaction(txObj)
}

// InvokeContract implements BlockchainClient interface
func (b *blockchainClientAdapter) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	return b.client.InvokeContract(contractHash, method, params)
}

// DeployContract implements BlockchainClient interface
func (b *blockchainClientAdapter) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	// Using empty signers and nil private key for now
	return b.client.DeployContract(ctx, nefFile, manifest, nil, nil)
}

// SubscribeToEvents implements BlockchainClient interface
func (b *blockchainClientAdapter) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	return b.client.SubscribeToEvents(ctx, contractHash, eventName, handler)
}

// GetTransactionReceipt implements BlockchainClient interface
func (b *blockchainClientAdapter) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	receipt, err := b.client.GetTransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// IsTransactionInMempool implements BlockchainClient interface
func (b *blockchainClientAdapter) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	return b.client.IsTransactionInMempool(ctx, hash)
}

// CheckHealth implements BlockchainClient interface
func (b *blockchainClientAdapter) CheckHealth(ctx context.Context) error {
	return b.client.CheckHealth(ctx)
}

// ResetConnections implements BlockchainClient interface
func (b *blockchainClientAdapter) ResetConnections() {
	b.client.ResetConnections()
}

// Close implements BlockchainClient interface
func (b *blockchainClientAdapter) Close() error {
	return b.client.Close()
}

// MockBlockchainClient is a mock implementation of BlockchainClient for testing
type MockBlockchainClient struct{}

// GetBlockHeight implements BlockchainClient interface for mock
func (m *MockBlockchainClient) GetBlockHeight() (uint32, error) {
	return 100, nil // Return mock height
}

// GetBlock implements BlockchainClient interface for mock
func (m *MockBlockchainClient) GetBlock(height uint32) (interface{}, error) {
	return map[string]interface{}{"height": height, "hash": "mock_hash"}, nil
}

// GetTransaction implements BlockchainClient interface for mock
func (m *MockBlockchainClient) GetTransaction(hash string) (interface{}, error) {
	return map[string]interface{}{"hash": hash, "status": "confirmed"}, nil
}

// SendTransaction implements BlockchainClient interface for mock
func (m *MockBlockchainClient) SendTransaction(tx interface{}) (string, error) {
	return "mock_tx_hash", nil
}

// InvokeContract implements BlockchainClient interface for mock
func (m *MockBlockchainClient) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"result": "mock_contract_result"}, nil
}

// DeployContract implements BlockchainClient interface for mock
func (m *MockBlockchainClient) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	return "mock_contract_hash", nil
}

// SubscribeToEvents implements BlockchainClient interface for mock
func (m *MockBlockchainClient) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	return nil // Mock does nothing
}

// GetTransactionReceipt implements BlockchainClient interface for mock
func (m *MockBlockchainClient) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	return map[string]interface{}{"hash": hash, "status": "success"}, nil
}

// IsTransactionInMempool implements BlockchainClient interface for mock
func (m *MockBlockchainClient) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	return true, nil // Always return true for testing
}

// CheckHealth implements BlockchainClient interface for mock
func (m *MockBlockchainClient) CheckHealth(ctx context.Context) error {
	return nil // Always healthy
}

// ResetConnections implements BlockchainClient interface for mock
func (m *MockBlockchainClient) ResetConnections() {
	// Mock does nothing
}

// Close implements BlockchainClient interface for mock
func (m *MockBlockchainClient) Close() error {
	return nil // Mock does nothing
}

// NewMockClient creates a new mock blockchain client for testing
func (f *ClientFactory) NewMockClient() BlockchainClient {
	return &MockBlockchainClient{}
}

// Initialize the global factory
func init() {
	// The factory will be initialized when the application starts
	// using the application's config and logger
}
