package blockchain

import (
	"context"
	"encoding/json"
	"github.com/R3E-Network/service_layer/pkg/logger"
)

// MockBlockchainClient implements the BlockchainClient interface for simple testing
type MockBlockchainClient struct {
	mockClient *MockClient
}

// NewMockBlockchainClient creates a new mock blockchain client
func NewMockBlockchainClient() *MockBlockchainClient {
	return &MockBlockchainClient{
		mockClient: NewMockClient(logger.NewLogger()),
	}
}

// GetBlockHeight implements BlockchainClient interface
func (m *MockBlockchainClient) GetBlockHeight() (uint32, error) {
	return m.mockClient.GetBlockHeight()
}

// GetBlock implements BlockchainClient interface
func (m *MockBlockchainClient) GetBlock(height uint32) (interface{}, error) {
	return m.mockClient.GetBlock(height)
}

// GetTransaction implements BlockchainClient interface
func (m *MockBlockchainClient) GetTransaction(hash string) (interface{}, error) {
	return m.mockClient.GetTransaction(hash)
}

// SendTransaction implements BlockchainClient interface
func (m *MockBlockchainClient) SendTransaction(tx interface{}) (string, error) {
	return m.mockClient.SendTransaction(tx)
}

// InvokeContract implements BlockchainClient interface
func (m *MockBlockchainClient) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	return m.mockClient.InvokeContract(contractHash, method, params)
}

// DeployContract implements BlockchainClient interface
func (m *MockBlockchainClient) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	return m.mockClient.DeployContract(ctx, nefFile, manifest)
}

// SubscribeToEvents implements BlockchainClient interface
func (m *MockBlockchainClient) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	return m.mockClient.SubscribeToEvents(ctx, contractHash, eventName, handler)
}

// GetTransactionReceipt implements BlockchainClient interface
func (m *MockBlockchainClient) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	return m.mockClient.GetTransactionReceipt(ctx, hash)
}

// IsTransactionInMempool implements BlockchainClient interface
func (m *MockBlockchainClient) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	return m.mockClient.IsTransactionInMempool(ctx, hash)
}

// CheckHealth implements BlockchainClient interface
func (m *MockBlockchainClient) CheckHealth(ctx context.Context) error {
	return m.mockClient.CheckHealth(ctx)
}

// ResetConnections implements BlockchainClient interface
func (m *MockBlockchainClient) ResetConnections() {
	m.mockClient.ResetConnections()
}

// Close implements BlockchainClient interface
func (m *MockBlockchainClient) Close() error {
	return m.mockClient.Close()
}
