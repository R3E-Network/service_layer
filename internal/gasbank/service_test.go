package gasbank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/neo-oracle/internal/config"
)

// Mock blockchain client for testing
type mockBlockchainClient struct{}

func (m *mockBlockchainClient) Connect() error {
	return nil
}

func (m *mockBlockchainClient) Disconnect() error {
	return nil
}

func (m *mockBlockchainClient) IsConnected() bool {
	return true
}

func (m *mockBlockchainClient) GetHeight() (uint32, error) {
	return 100, nil
}

func (m *mockBlockchainClient) SwitchEndpoint() error {
	return nil
}

func TestNewService(t *testing.T) {
	cfg := &config.GasBankConfig{
		MinimumGasBalance: 10.0,
		AutoRefill:        true,
		RefillAmount:      50.0,
	}

	bc := &mockBlockchainClient{}
	svc := NewService(cfg, bc)

	assert.NotNil(t, svc)
	assert.Equal(t, cfg, svc.cfg)
	assert.Equal(t, bc, svc.blockchain)
}

func TestGasOperations(t *testing.T) {
	cfg := &config.GasBankConfig{
		MinimumGasBalance: 10.0,
		AutoRefill:        true,
		RefillAmount:      50.0,
	}

	bc := &mockBlockchainClient{}
	svc := NewService(cfg, bc)

	// Start the service
	ctx := context.Background()
	err := svc.Start(ctx)
	assert.NoError(t, err)

	// Test deposit
	userID := "test-user"
	amount := 100.0
	err = svc.DepositGas(userID, amount)
	assert.NoError(t, err)

	// Test get balance
	balance, err := svc.GetBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, amount, balance)

	// Test withdraw
	withdrawAmount := 25.0
	err = svc.WithdrawGas(userID, withdrawAmount)
	assert.NoError(t, err)

	// Check updated balance
	balance, err = svc.GetBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, amount-withdrawAmount, balance)

	// Test allocation
	functionID := "test-function"
	allocAmount := 10.0
	allocID, err := svc.AllocateGas(userID, functionID, allocAmount)
	assert.NoError(t, err)
	assert.NotEmpty(t, allocID)

	// Check balance after allocation
	balance, err = svc.GetBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, amount-withdrawAmount-allocAmount, balance)

	// Test finalize
	actualAmount := 5.0
	err = svc.FinalizeGasUsage(allocID, actualAmount)
	assert.NoError(t, err)

	// Stop the service
	err = svc.Stop()
	assert.NoError(t, err)
}

func TestInvalidOperations(t *testing.T) {
	cfg := &config.GasBankConfig{
		MinimumGasBalance: 10.0,
		AutoRefill:        true,
		RefillAmount:      50.0,
	}

	bc := &mockBlockchainClient{}
	svc := NewService(cfg, bc)

	// Test deposit with negative amount
	userID := "test-user"
	err := svc.DepositGas(userID, -10.0)
	assert.Error(t, err)

	// Test withdraw with no balance
	err = svc.WithdrawGas(userID, 10.0)
	assert.Error(t, err)

	// Test allocate with no balance
	_, err = svc.AllocateGas(userID, "test-function", 10.0)
	assert.Error(t, err)

	// Test unknown user
	_, err = svc.GetBalance("unknown-user")
	assert.Error(t, err)
}
