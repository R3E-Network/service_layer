package gasbank

import (
	"testing"
	"time"

	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CoreGasBankServiceInterface defines the interface for the core gasbank service used in tests
type CoreGasBankServiceInterface interface {
	CreateAccount(userID string, walletAddress string) (*models.GasBankAccount, error)
	GetAccount(id string) (*models.GasBankAccount, error)
	GetAccountByUserID(userID string) (*models.GasBankAccount, error)
	GetAccountByWalletAddress(walletAddress string) (*models.GasBankAccount, error)
	ListAccounts() ([]*models.GasBankAccount, error)
	GetBalance(accountID string) (string, error)
	GetAvailableBalance(accountID string) (string, error)
	ProcessDeposit(fromAddress string, toAddress string, amount string, blockchainTxID string, blockHeight uint32) (*models.GasBankTransaction, error)
	RequestWithdrawal(userID string, amount string, toAddress string) (*models.WithdrawalRequest, error)
	ProcessWithdrawalRequest(requestID string) (*models.GasBankTransaction, error)
	CancelWithdrawalRequest(requestID string) error
	DeductFee(userID string, amount string, notes string) (*models.GasBankTransaction, error)
	GetTransaction(id string) (*models.GasBankTransaction, error)
	ListTransactionsByUserID(userID string, limit int, offset int) ([]*models.GasBankTransaction, error)
	ListTransactionsByAccountID(accountID string, limit int, offset int) ([]*models.GasBankTransaction, error)
	GetWithdrawalRequest(id string) (*models.WithdrawalRequest, error)
	ListWithdrawalRequestsByUserID(userID string, limit int, offset int) ([]*models.WithdrawalRequest, error)
	Start() error
	Stop()
}

// MockCoreGasBankService is a mock implementation of the core GasBank service
type MockCoreGasBankService struct {
	mock.Mock
}

// Implement all methods required by the core GasBank service

func (m *MockCoreGasBankService) CreateAccount(userID string, walletAddress string) (*models.GasBankAccount, error) {
	args := m.Called(userID, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankAccount), args.Error(1)
}

func (m *MockCoreGasBankService) GetAccount(id string) (*models.GasBankAccount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankAccount), args.Error(1)
}

func (m *MockCoreGasBankService) GetAccountByUserID(userID string) (*models.GasBankAccount, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankAccount), args.Error(1)
}

func (m *MockCoreGasBankService) GetAccountByWalletAddress(walletAddress string) (*models.GasBankAccount, error) {
	args := m.Called(walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankAccount), args.Error(1)
}

func (m *MockCoreGasBankService) ListAccounts() ([]*models.GasBankAccount, error) {
	args := m.Called()
	return args.Get(0).([]*models.GasBankAccount), args.Error(1)
}

func (m *MockCoreGasBankService) GetBalance(accountID string) (string, error) {
	args := m.Called(accountID)
	return args.String(0), args.Error(1)
}

func (m *MockCoreGasBankService) GetAvailableBalance(accountID string) (string, error) {
	args := m.Called(accountID)
	return args.String(0), args.Error(1)
}

func (m *MockCoreGasBankService) ProcessDeposit(fromAddress string, toAddress string, amount string, blockchainTxID string, blockHeight uint32) (*models.GasBankTransaction, error) {
	args := m.Called(fromAddress, toAddress, amount, blockchainTxID, blockHeight)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankTransaction), args.Error(1)
}

func (m *MockCoreGasBankService) RequestWithdrawal(userID string, amount string, toAddress string) (*models.WithdrawalRequest, error) {
	args := m.Called(userID, amount, toAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WithdrawalRequest), args.Error(1)
}

func (m *MockCoreGasBankService) ProcessWithdrawalRequest(requestID string) (*models.GasBankTransaction, error) {
	args := m.Called(requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankTransaction), args.Error(1)
}

func (m *MockCoreGasBankService) CancelWithdrawalRequest(requestID string) error {
	args := m.Called(requestID)
	return args.Error(0)
}

func (m *MockCoreGasBankService) DeductFee(userID string, amount string, notes string) (*models.GasBankTransaction, error) {
	args := m.Called(userID, amount, notes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankTransaction), args.Error(1)
}

func (m *MockCoreGasBankService) GetTransaction(id string) (*models.GasBankTransaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GasBankTransaction), args.Error(1)
}

func (m *MockCoreGasBankService) ListTransactionsByUserID(userID string, limit int, offset int) ([]*models.GasBankTransaction, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*models.GasBankTransaction), args.Error(1)
}

func (m *MockCoreGasBankService) ListTransactionsByAccountID(accountID string, limit int, offset int) ([]*models.GasBankTransaction, error) {
	args := m.Called(accountID, limit, offset)
	return args.Get(0).([]*models.GasBankTransaction), args.Error(1)
}

func (m *MockCoreGasBankService) GetWithdrawalRequest(id string) (*models.WithdrawalRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WithdrawalRequest), args.Error(1)
}

func (m *MockCoreGasBankService) ListWithdrawalRequestsByUserID(userID string, limit int, offset int) ([]*models.WithdrawalRequest, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*models.WithdrawalRequest), args.Error(1)
}

func (m *MockCoreGasBankService) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCoreGasBankService) Stop() {
	m.Called()
}

// TestWrapper is a wrapper around the mock core GasBank service for testing
type TestWrapper struct {
	coreService CoreGasBankServiceInterface
}

// Test helper functions
func setupWrapperTest(t *testing.T) (*TestWrapper, *MockCoreGasBankService) {
	mockCoreService := new(MockCoreGasBankService)
	// Create a test wrapper with the mock service
	wrapper := &TestWrapper{
		coreService: mockCoreService,
	}
	return wrapper, mockCoreService
}

// Test cases
func TestWrapperCreateAccount(t *testing.T) {
	wrapper, mockCoreService := setupWrapperTest(t)

	// Setup test parameters
	userID := "user123"
	walletAddress := "0xabcdef1234567890"

	// Setup expected result
	expectedAccount := &models.GasBankAccount{
		ID:            "acc123",
		UserID:        userID,
		WalletAddress: walletAddress,
		Balance:       "100.0",
	}

	// Setup mock expectations
	mockCoreService.On("CreateAccount", userID, walletAddress).Return(expectedAccount, nil)

	// Call the wrapper method
	result, err := wrapper.coreService.CreateAccount(userID, walletAddress)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, result)
	mockCoreService.AssertExpectations(t)
}

func TestWrapperGetBalance(t *testing.T) {
	wrapper, mockCoreService := setupWrapperTest(t)

	// Setup test parameters
	accountID := "acc123"
	expectedBalance := "500.75"

	// Setup mock expectations
	mockCoreService.On("GetBalance", accountID).Return(expectedBalance, nil)

	// Call the wrapper method
	balance, err := wrapper.coreService.GetBalance(accountID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	mockCoreService.AssertExpectations(t)
}

func TestWrapperRequestWithdrawal(t *testing.T) {
	wrapper, mockCoreService := setupWrapperTest(t)

	// Setup test parameters
	userID := "user123"
	amount := "50.0"
	toAddress := "0xdef456"

	// Setup expected result
	createdTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	expectedRequest := &models.WithdrawalRequest{
		ID:        "w123",
		UserID:    userID,
		AccountID: "acc123",
		Amount:    amount,
		ToAddress: toAddress,
		Status:    "pending",
		CreatedAt: createdTime,
	}

	// Setup mock expectations
	mockCoreService.On("RequestWithdrawal", userID, amount, toAddress).Return(expectedRequest, nil)

	// Call the wrapper method
	result, err := wrapper.coreService.RequestWithdrawal(userID, amount, toAddress)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, result)
	mockCoreService.AssertExpectations(t)
}

func TestWrapperProcessDeposit(t *testing.T) {
	wrapper, mockCoreService := setupWrapperTest(t)

	// Setup test parameters
	fromAddress := "0xsender"
	toAddress := "0xreceiver"
	amount := "100.0"
	blockchainTxID := "tx123"
	blockHeight := uint32(12345)

	// Setup expected result
	expectedTransaction := &models.GasBankTransaction{
		ID:             "tr123",
		AccountID:      "acc123",
		Type:           models.DepositTransaction,
		Amount:         amount,
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		BlockchainTxID: blockchainTxID,
		BlockHeight:    blockHeight,
	}

	// Setup mock expectations
	mockCoreService.On("ProcessDeposit", fromAddress, toAddress, amount, blockchainTxID, blockHeight).Return(expectedTransaction, nil)

	// Call the wrapper method
	result, err := wrapper.coreService.ProcessDeposit(fromAddress, toAddress, amount, blockchainTxID, blockHeight)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedTransaction, result)
	mockCoreService.AssertExpectations(t)
}

func TestWrapperStartAndStop(t *testing.T) {
	wrapper, mockCoreService := setupWrapperTest(t)

	// Setup mock expectations
	mockCoreService.On("Start").Return(nil)
	mockCoreService.On("Stop").Return()

	// Call the wrapper methods
	err := wrapper.coreService.Start()
	assert.NoError(t, err)
	
	// Call Stop and check expectations
	wrapper.coreService.Stop()
	mockCoreService.AssertExpectations(t)
}