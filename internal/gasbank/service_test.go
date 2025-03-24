package gasbank

import (
	"context"
	"testing"
	"time"

	"github.com/R3E-Network/service_layer/internal/blockchain"
	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/R3E-Network/service_layer/internal/tee"
	"github.com/R3E-Network/service_layer/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock blockchain client for testing
type mockBlockchainClient struct{}

// Connect implements the Client interface
func (m *mockBlockchainClient) Connect() error {
	return nil
}

// Disconnect implements the Client interface
func (m *mockBlockchainClient) Disconnect() error {
	return nil
}

// IsConnected implements the Client interface
func (m *mockBlockchainClient) IsConnected() bool {
	return true
}

// GetHeight implements the Client interface
func (m *mockBlockchainClient) GetHeight() (uint32, error) {
	return 100, nil
}

// SwitchEndpoint implements the Client interface
func (m *mockBlockchainClient) SwitchEndpoint() error {
	return nil
}

// CreateTransaction implements the Client interface with the correct signature
func (m *mockBlockchainClient) CreateTransaction(ctx context.Context, params blockchain.TransactionParams) (string, error) {
	return "tx-12345", nil
}

// Additional methods needed for the blockchain.Client interface
func (m *mockBlockchainClient) GetBlockCount(ctx context.Context) (int, error) {
	return 100, nil
}

func (m *mockBlockchainClient) InvokeFunction(ctx context.Context, scriptHash, operation string, params []interface{}) (blockchain.InvocationResult, error) {
	return blockchain.InvocationResult{State: "HALT"}, nil
}

func (m *mockBlockchainClient) SignTransaction(ctx context.Context, tx string, privateKey string) (string, error) {
	return "signed-tx", nil
}

func (m *mockBlockchainClient) SendTransaction(ctx context.Context, signedTx string) (string, error) {
	return "tx-hash", nil
}

func (m *mockBlockchainClient) GetTransaction(ctx context.Context, txid string) (blockchain.Transaction, error) {
	return blockchain.Transaction{ID: txid}, nil
}

func (m *mockBlockchainClient) GetStorage(ctx context.Context, scriptHash string, key string) (string, error) {
	return "storage-value", nil
}

func (m *mockBlockchainClient) GetBalance(ctx context.Context, address string, assetID string) (string, error) {
	return "100.0", nil
}

// Mock repository for testing
type mockGasBankRepository struct {
	accounts     map[string]*models.GasBankAccount
	transactions map[string]*models.GasBankTransaction
	withdrawals  map[string]*models.WithdrawalRequest
	deposits     map[string]*models.DepositTracker
}

// NewMockGasBankRepository creates a new mock repository
func NewMockGasBankRepository() *mockGasBankRepository {
	userID := "test-user"
	testAccount := &models.GasBankAccount{
		ID:               "acc-12345",
		UserID:           userID,
		WalletAddress:    "wallet-12345",
		Balance:          "100.0",
		AvailableBalance: "100.0",
		PendingBalance:   "0.0",
		Active:           true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	accounts := map[string]*models.GasBankAccount{
		testAccount.ID: testAccount,
	}

	return &mockGasBankRepository{
		accounts:     accounts,
		transactions: make(map[string]*models.GasBankTransaction),
		withdrawals:  make(map[string]*models.WithdrawalRequest),
		deposits:     make(map[string]*models.DepositTracker),
	}
}

// Account operations
func (m *mockGasBankRepository) CreateAccount(ctx interface{}, account *models.GasBankAccount) (*models.GasBankAccount, error) {
	m.accounts[account.ID] = account
	return account, nil
}

func (m *mockGasBankRepository) GetAccount(ctx interface{}, id string) (*models.GasBankAccount, error) {
	account, exists := m.accounts[id]
	if !exists {
		return nil, nil
	}
	return account, nil
}

func (m *mockGasBankRepository) GetAccountByUserID(ctx interface{}, userID string) (*models.GasBankAccount, error) {
	for _, account := range m.accounts {
		if account.UserID == userID {
			return account, nil
		}
	}
	return nil, nil
}

func (m *mockGasBankRepository) GetAccountByWalletAddress(ctx interface{}, address string) (*models.GasBankAccount, error) {
	for _, account := range m.accounts {
		if account.WalletAddress == address {
			return account, nil
		}
	}
	return nil, nil
}

func (m *mockGasBankRepository) UpdateAccount(ctx interface{}, account *models.GasBankAccount) (*models.GasBankAccount, error) {
	m.accounts[account.ID] = account
	return account, nil
}

func (m *mockGasBankRepository) ListAccounts(ctx interface{}) ([]*models.GasBankAccount, error) {
	accounts := make([]*models.GasBankAccount, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// Transaction operations
func (m *mockGasBankRepository) CreateTransaction(ctx interface{}, tx *models.GasBankTransaction) (*models.GasBankTransaction, error) {
	m.transactions[tx.ID] = tx
	return tx, nil
}

func (m *mockGasBankRepository) GetTransaction(ctx interface{}, id string) (*models.GasBankTransaction, error) {
	tx, exists := m.transactions[id]
	if !exists {
		return nil, nil
	}
	return tx, nil
}

func (m *mockGasBankRepository) GetTransactionByBlockchainTxID(ctx interface{}, txID string) (*models.GasBankTransaction, error) {
	for _, tx := range m.transactions {
		if tx.BlockchainTxID == txID {
			return tx, nil
		}
	}
	return nil, nil
}

func (m *mockGasBankRepository) UpdateTransaction(ctx interface{}, tx *models.GasBankTransaction) (*models.GasBankTransaction, error) {
	m.transactions[tx.ID] = tx
	return tx, nil
}

func (m *mockGasBankRepository) ListTransactionsByUserID(ctx interface{}, userID string, limit int, offset int) ([]*models.GasBankTransaction, error) {
	var result []*models.GasBankTransaction
	for _, tx := range m.transactions {
		if tx.UserID == userID {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockGasBankRepository) ListTransactionsByAccountID(ctx interface{}, accountID string, limit int, offset int) ([]*models.GasBankTransaction, error) {
	var result []*models.GasBankTransaction
	for _, tx := range m.transactions {
		if tx.AccountID == accountID {
			result = append(result, tx)
		}
	}
	return result, nil
}

// Withdrawal operations
func (m *mockGasBankRepository) CreateWithdrawalRequest(ctx interface{}, req *models.WithdrawalRequest) (*models.WithdrawalRequest, error) {
	m.withdrawals[req.ID] = req
	return req, nil
}

func (m *mockGasBankRepository) GetWithdrawalRequest(ctx interface{}, id string) (*models.WithdrawalRequest, error) {
	req, exists := m.withdrawals[id]
	if !exists {
		return nil, nil
	}
	return req, nil
}

func (m *mockGasBankRepository) UpdateWithdrawalRequest(ctx interface{}, req *models.WithdrawalRequest) (*models.WithdrawalRequest, error) {
	m.withdrawals[req.ID] = req
	return req, nil
}

func (m *mockGasBankRepository) ListWithdrawalRequestsByUserID(ctx interface{}, userID string, limit int, offset int) ([]*models.WithdrawalRequest, error) {
	var result []*models.WithdrawalRequest
	for _, req := range m.withdrawals {
		if req.UserID == userID {
			result = append(result, req)
		}
	}
	return result, nil
}

// Deposit tracking operations
func (m *mockGasBankRepository) CreateDepositTracker(ctx interface{}, deposit *models.DepositTracker) (*models.DepositTracker, error) {
	m.deposits[deposit.ID] = deposit
	return deposit, nil
}

func (m *mockGasBankRepository) GetDepositTrackerByTxID(ctx interface{}, txID string) (*models.DepositTracker, error) {
	for _, deposit := range m.deposits {
		if deposit.BlockchainTxID == txID {
			return deposit, nil
		}
	}
	return nil, nil
}

func (m *mockGasBankRepository) UpdateDepositTracker(ctx interface{}, deposit *models.DepositTracker) (*models.DepositTracker, error) {
	m.deposits[deposit.ID] = deposit
	return deposit, nil
}

func (m *mockGasBankRepository) ListUnprocessedDeposits(ctx interface{}) ([]*models.DepositTracker, error) {
	var result []*models.DepositTracker
	for _, deposit := range m.deposits {
		if !deposit.Processed {
			result = append(result, deposit)
		}
	}
	return result, nil
}

// Balance operations
func (m *mockGasBankRepository) UpdateBalance(ctx interface{}, accountID string, newBalance string, newPendingBalance string, newAvailableBalance string) error {
	account, exists := m.accounts[accountID]
	if !exists {
		return nil
	}
	account.Balance = newBalance
	account.PendingBalance = newPendingBalance
	account.AvailableBalance = newAvailableBalance
	return nil
}

func (m *mockGasBankRepository) IncrementDailyWithdrawal(ctx interface{}, accountID string, amount string) error {
	return nil
}

func (m *mockGasBankRepository) ResetDailyWithdrawal(ctx interface{}, accountID string) error {
	return nil
}

func TestNewService(t *testing.T) {
	// Create Config with GasBankConfig
	cfg := &config.Config{
		GasBank: config.GasBankConfig{
			MinimumGasBalance: 10.0,
			AutoRefill:        true,
			RefillAmount:      50.0,
		},
	}

	// Create logger
	log := logger.New(logger.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	})

	// Create TEE manager
	teeManager, err := tee.NewManager(&config.Config{
		TEE: config.TEEConfig{
			Provider:          "simulation",
			EnableAttestation: false,
		},
		Functions: config.FunctionsConfig{
			MaxMemory:        512,
			ExecutionTimeout: 30,
			MaxConcurrency:   10,
		},
	}, log)
	require.NoError(t, err)

	// Create mock dependencies
	repository := NewMockGasBankRepository()
	bc := &mockBlockchainClient{}
	
	// Create service
	svc, err := NewService(cfg, repository, bc, teeManager)
	require.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestGasOperations(t *testing.T) {
	// Create Config with GasBankConfig
	cfg := &config.Config{
		GasBank: config.GasBankConfig{
			MinimumGasBalance: 10.0,
			AutoRefill:        true,
			RefillAmount:      50.0,
		},
	}

	// Create logger
	log := logger.New(logger.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	})

	// Create TEE manager
	teeManager, err := tee.NewManager(&config.Config{
		TEE: config.TEEConfig{
			Provider:          "simulation",
			EnableAttestation: false,
		},
		Functions: config.FunctionsConfig{
			MaxMemory:        512,
			ExecutionTimeout: 30,
			MaxConcurrency:   10,
		},
	}, log)
	require.NoError(t, err)

	// Create mock dependencies
	repository := NewMockGasBankRepository()
	bc := &mockBlockchainClient{}
	
	// Create service
	svc, err := NewService(cfg, repository, bc, teeManager)
	require.NoError(t, err)

	// Start the service
	ctx := context.Background()
	err = svc.Start(ctx)
	assert.NoError(t, err)

	// Create account
	userID := "test-user"
	walletAddress := "wallet-12345"
	account, err := svc.CreateAccount(ctx, userID, walletAddress)
	assert.NoError(t, err)
	assert.NotNil(t, account)

	// Test get account
	retrievedAccount, err := svc.GetAccountByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, account.ID, retrievedAccount.ID)

	// Test process deposit
	fromAddress := "deposit-source"
	toAddress := walletAddress
	amount := "100.0"
	blockchainTxID := "tx-12345"
	deposit, err := svc.ProcessDeposit(ctx, fromAddress, toAddress, amount, blockchainTxID, 100)
	assert.NoError(t, err)
	assert.NotNil(t, deposit)

	// Test get balance
	balance, err := svc.GetBalance(ctx, account.ID)
	assert.NoError(t, err)
	assert.Equal(t, "100.0", balance)

	// Test request withdrawal
	withdrawAmount := "25.0"
	withdrawalRequest, err := svc.RequestWithdrawal(ctx, userID, withdrawAmount, "withdraw-dest")
	assert.NoError(t, err)
	assert.NotNil(t, withdrawalRequest)

	// Test process withdrawal request
	withdrawal, err := svc.ProcessWithdrawalRequest(ctx, withdrawalRequest.ID)
	assert.NoError(t, err)
	assert.NotNil(t, withdrawal)

	// Test deduct fee
	feeAmount := "10.0"
	fee, err := svc.DeductFee(ctx, userID, feeAmount, "Test fee")
	assert.NoError(t, err)
	assert.NotNil(t, fee)

	// Stop the service
	err = svc.Stop(ctx)
	assert.NoError(t, err)
}

func TestInvalidOperations(t *testing.T) {
	// Create Config with GasBankConfig
	cfg := &config.Config{
		GasBank: config.GasBankConfig{
			MinimumGasBalance: 10.0,
			AutoRefill:        true,
			RefillAmount:      50.0,
		},
	}

	// Create logger
	log := logger.New(logger.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	})

	// Create TEE manager
	teeManager, err := tee.NewManager(&config.Config{
		TEE: config.TEEConfig{
			Provider:          "simulation",
			EnableAttestation: false,
		},
		Functions: config.FunctionsConfig{
			MaxMemory:        512,
			ExecutionTimeout: 30,
			MaxConcurrency:   10,
		},
	}, log)
	require.NoError(t, err)

	// Create mock dependencies
	repository := NewMockGasBankRepository()
	bc := &mockBlockchainClient{}
	
	// Create service
	svc, err := NewService(cfg, repository, bc, teeManager)
	require.NoError(t, err)

	// Test with non-existent user
	ctx := context.Background()
	_, err = svc.GetAccountByUserID(ctx, "non-existent-user")
	assert.Error(t, err)

	// Test withdrawal with insufficient balance
	userID := "test-user"
	walletAddress := "wallet-12345"
	_, err = svc.CreateAccount(ctx, userID, walletAddress)
	assert.NoError(t, err)
	
	// Try to withdraw more than available
	_, err = svc.RequestWithdrawal(ctx, userID, "1000.0", "withdraw-dest")
	assert.Error(t, err)
}
