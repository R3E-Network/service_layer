package blockchain

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBlockchainClientForTests implements the BlockchainClient interface for testing
type MockBlockchainClientForTests struct {
	mock.Mock
}

// GetBlockHeight mocks the GetBlockHeight method
func (m *MockBlockchainClientForTests) GetBlockHeight() (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}

// GetBlock mocks the GetBlock method
func (m *MockBlockchainClientForTests) GetBlock(height uint32) (interface{}, error) {
	args := m.Called(height)
	return args.Get(0), args.Error(1)
}

// GetTransaction mocks the GetTransaction method
func (m *MockBlockchainClientForTests) GetTransaction(hash string) (interface{}, error) {
	args := m.Called(hash)
	return args.Get(0), args.Error(1)
}

// SendTransaction mocks the SendTransaction method
func (m *MockBlockchainClientForTests) SendTransaction(tx interface{}) (string, error) {
	args := m.Called(tx)
	return args.String(0), args.Error(1)
}

// InvokeContract mocks the InvokeContract method
func (m *MockBlockchainClientForTests) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	args := m.Called(contractHash, method, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// DeployContract mocks the DeployContract method
func (m *MockBlockchainClientForTests) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	args := m.Called(ctx, nefFile, manifest)
	return args.String(0), args.Error(1)
}

// SubscribeToEvents mocks the SubscribeToEvents method
func (m *MockBlockchainClientForTests) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	args := m.Called(ctx, contractHash, eventName, handler)
	return args.Error(0)
}

// GetTransactionReceipt mocks the GetTransactionReceipt method
func (m *MockBlockchainClientForTests) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	args := m.Called(ctx, hash)
	return args.Get(0), args.Error(1)
}

// IsTransactionInMempool mocks the IsTransactionInMempool method
func (m *MockBlockchainClientForTests) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	args := m.Called(ctx, hash)
	return args.Bool(0), args.Error(1)
}

// CheckHealth mocks the CheckHealth method
func (m *MockBlockchainClientForTests) CheckHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ResetConnections mocks the ResetConnections method
func (m *MockBlockchainClientForTests) ResetConnections() {
	m.Called()
}

// Close mocks the Close method
func (m *MockBlockchainClientForTests) Close() error {
	args := m.Called()
	return args.Error(0)
}

// CreateMockBlockchainClient creates a new MockBlockchainClientForTests
func CreateMockBlockchainClient() *MockBlockchainClientForTests {
	return &MockBlockchainClientForTests{}
}

// SetupTestBlockchain creates a mock blockchain client for testing
func SetupTestBlockchain(t *testing.T) (BlockchainClient, func()) {
	t.Helper()

	// Create a mock client
	mockClient := CreateMockBlockchainClient()

	// Return the mock client and a cleanup function
	return mockClient, func() {
		// Cleanup (close the client)
		err := mockClient.Close()
		require.NoError(t, err)
	}
}

// SetupTestFactory creates a blockchain client factory for testing
func SetupTestFactory(t *testing.T) BlockchainClientFactory {
	t.Helper()

	// Create a minimal config for testing
	neoConfig := &config.NeoConfig{
		NetworkID:   860833102,
		ChainID:     1,
		RPCEndpoint: "http://localhost:10332",
	}

	// Create a factory with nil logger for testing
	return NewNeoBlockchainClientFactory(neoConfig, nil)
}

// SetupTestContract sets up a mock smart contract in the blockchain client
func SetupTestContract(t *testing.T, client *MockBlockchainClientForTests, contractHash string, methods map[string]interface{}) {
	t.Helper()

	// Add contract methods and their results
	for method, result := range methods {
		resultMap := map[string]interface{}{
			"result": result,
		}
		client.On("InvokeContract", contractHash, method, []interface{}{}).Return(resultMap, nil)
	}
}

// CreateDummyContractManifest creates a dummy contract manifest for testing
func CreateDummyContractManifest(name, description string) json.RawMessage {
	manifest := map[string]interface{}{
		"name":        name,
		"description": description,
		"abi": map[string]interface{}{
			"methods": []map[string]interface{}{
				{
					"name":       "transfer",
					"parameters": []string{"Hash160", "Hash160", "Integer"},
					"returnType": "Boolean",
				},
				{
					"name":       "balanceOf",
					"parameters": []string{"Hash160"},
					"returnType": "Integer",
				},
			},
		},
	}

	data, _ := json.Marshal(manifest)
	return data
}
