package mocks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/R3E-Network/service_layer/internal/blockchain"
	"github.com/R3E-Network/service_layer/internal/database"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository implements the database.TransactionRepository interface for testing
type MockTransactionRepository struct {
	mock.Mock
}

// CreateTransaction mocks the CreateTransaction method
func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

// GetTransactionByID mocks the GetTransactionByID method
func (m *MockTransactionRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

// ListTransactions mocks the ListTransactions method
func (m *MockTransactionRepository) ListTransactions(
	ctx context.Context,
	service string,
	status models.TransactionStatus,
	entityID *uuid.UUID,
	page, limit int,
) (*models.TransactionListResponse, error) {
	args := m.Called(ctx, service, status, entityID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	// Create a proper TransactionListResponse
	txns := args.Get(0).([]models.Transaction)
	total := args.Get(1).(int)
	
	return &models.TransactionListResponse{
		Transactions: txns,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, args.Error(2)
}

// UpdateTransactionStatus mocks the UpdateTransactionStatus method
func (m *MockTransactionRepository) UpdateTransactionStatus(
	ctx context.Context,
	id uuid.UUID,
	status models.TransactionStatus,
	result json.RawMessage,
	gasConsumed *int64,
	blockHeight *int64,
	blockTime *time.Time,
	errorMessage string,
) error {
	args := m.Called(ctx, id, status, result, gasConsumed, blockHeight, blockTime, errorMessage)
	return args.Error(0)
}

// GetWalletByService mocks the GetWalletByService method
func (m *MockTransactionRepository) GetWalletByService(ctx context.Context, service string) (*models.WalletAccount, error) {
	args := m.Called(ctx, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WalletAccount), args.Error(1)
}

// GetWalletByAddress retrieves a wallet account by its address
func (m *MockTransactionRepository) GetWalletByAddress(ctx context.Context, address string) (*models.WalletAccount, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WalletAccount), args.Error(1)
}

// AddTransactionEvent mocks the AddTransactionEvent method
func (m *MockTransactionRepository) AddTransactionEvent(ctx context.Context, event *models.TransactionEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// DeleteTransaction mocks the DeleteTransaction method
func (m *MockTransactionRepository) DeleteTransaction(ctx context.Context, txID uuid.UUID) error {
	args := m.Called(ctx, txID)
	return args.Error(0)
}

// DeleteWalletAccount mocks the DeleteWalletAccount method
func (m *MockTransactionRepository) DeleteWalletAccount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CreateWalletAccount mocks the CreateWalletAccount method
func (m *MockTransactionRepository) CreateWalletAccount(ctx context.Context, wallet *models.WalletAccount) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

// GetTransactionEvents mocks the GetTransactionEvents method
func (m *MockTransactionRepository) GetTransactionEvents(ctx context.Context, txID uuid.UUID) ([]models.TransactionEvent, error) {
	args := m.Called(ctx, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TransactionEvent), args.Error(1)
}

// GetTransactionByHash mocks the GetTransactionByHash method
func (m *MockTransactionRepository) GetTransactionByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

// ListWalletsByService mocks the ListWalletsByService method
func (m *MockTransactionRepository) ListWalletsByService(ctx context.Context, service string) ([]models.WalletAccount, error) {
	args := m.Called(ctx, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.WalletAccount), args.Error(1)
}

// MockClient implements a client interface for testing
type MockClient struct {
	mock.Mock
}

// GetBalance mocks the GetBalance method
func (m *MockClient) GetBalance(ctx context.Context, address string, assetID string) (string, error) {
	args := m.Called(ctx, address, assetID)
	return args.String(0), args.Error(1)
}

// CreateTransaction mocks the CreateTransaction method
func (m *MockClient) CreateTransaction(ctx context.Context, params blockchain.TransactionParams) (string, error) {
	args := m.Called(ctx, params)
	return args.String(0), args.Error(1)
}

// GetHeight mocks the GetHeight method
func (m *MockClient) GetHeight() (uint32, error) {
	args := m.Called()
	return uint32(args.Int(0)), args.Error(1)
}

// GetBlockCount mocks the GetBlockCount method
func (m *MockClient) GetBlockCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// SignTransaction mocks the SignTransaction method
func (m *MockClient) SignTransaction(ctx context.Context, tx string, privateKey string) (string, error) {
	args := m.Called(ctx, tx, privateKey)
	return args.String(0), args.Error(1)
}

// SendTransaction mocks the SendTransaction method
func (m *MockClient) SendTransaction(ctx context.Context, signedTx string) (string, error) {
	args := m.Called(ctx, signedTx)
	return args.String(0), args.Error(1)
}

// GetTransaction mocks the GetTransaction method
func (m *MockClient) GetTransaction(ctx context.Context, txid string) (blockchain.Transaction, error) {
	args := m.Called(ctx, txid)
	return args.Get(0).(blockchain.Transaction), args.Error(1)
}

// GetStorage mocks the GetStorage method
func (m *MockClient) GetStorage(ctx context.Context, scriptHash string, key string) (string, error) {
	args := m.Called(ctx, scriptHash, key)
	return args.String(0), args.Error(1)
}

// InvokeFunction mocks the InvokeFunction method
func (m *MockClient) InvokeFunction(ctx context.Context, scriptHash, operation string, params []interface{}) (blockchain.InvocationResult, error) {
	args := m.Called(ctx, scriptHash, operation, params)
	return args.Get(0).(blockchain.InvocationResult), args.Error(1)
}

// CheckHealth mocks the CheckHealth method
func (m *MockClient) CheckHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockWalletStore implements the blockchain.WalletStore interface for testing
type MockWalletStore struct {
	mock.Mock
}

// GetWallet mocks the GetWallet method
func (m *MockWalletStore) GetWallet(userID int, walletID string) (*wallet.Wallet, error) {
	args := m.Called(userID, walletID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Wallet), args.Error(1)
}

// MockBlockchainClient implements the blockchain.BlockchainClient interface for testing
type MockBlockchainClient struct {
	mock.Mock
}

// GetBlockHeight mocks the GetBlockHeight method
func (m *MockBlockchainClient) GetBlockHeight() (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}

// GetBlock mocks the GetBlock method
func (m *MockBlockchainClient) GetBlock(height uint32) (interface{}, error) {
	args := m.Called(height)
	return args.Get(0), args.Error(1)
}

// GetTransaction mocks the GetTransaction method
func (m *MockBlockchainClient) GetTransaction(hash string) (interface{}, error) {
	args := m.Called(hash)
	return args.Get(0), args.Error(1)
}

// SendTransaction mocks the SendTransaction method
func (m *MockBlockchainClient) SendTransaction(tx interface{}) (string, error) {
	args := m.Called(tx)
	return args.String(0), args.Error(1)
}

// InvokeContract mocks the InvokeContract method
func (m *MockBlockchainClient) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	args := m.Called(contractHash, method, params)
	if ret := args.Get(0); ret != nil {
		return ret.(map[string]interface{}), args.Error(1)
	}
	return nil, args.Error(1)
}

// DeployContract mocks the DeployContract method
func (m *MockBlockchainClient) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	args := m.Called(ctx, nefFile, manifest)
	return args.String(0), args.Error(1)
}

// SubscribeToEvents mocks the SubscribeToEvents method
func (m *MockBlockchainClient) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	args := m.Called(ctx, contractHash, eventName, handler)
	return args.Error(0)
}

// GetTransactionReceipt mocks the GetTransactionReceipt method
func (m *MockBlockchainClient) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	args := m.Called(ctx, hash)
	return args.Get(0), args.Error(1)
}

// IsTransactionInMempool mocks the IsTransactionInMempool method
func (m *MockBlockchainClient) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	args := m.Called(ctx, hash)
	return args.Bool(0), args.Error(1)
}

// CheckHealth mocks the CheckHealth method
func (m *MockBlockchainClient) CheckHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ResetConnections mocks the ResetConnections method
func (m *MockBlockchainClient) ResetConnections() {
	m.Called()
}

// Close mocks the Close method
func (m *MockBlockchainClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// CreateMockBlockchainClient creates a new MockBlockchainClient
func CreateMockBlockchainClient() *MockBlockchainClient {
	return &MockBlockchainClient{}
}

// Helper function to create mock services and a transaction service for testing
func CreateMockServices() (*MockTransactionRepository, *MockClient, *MockWalletStore, *blockchain.TransactionService) {
	mockRepo := new(MockTransactionRepository)
	mockClient := new(MockClient)
	mockWalletStore := new(MockWalletStore)

	// Set up any default expectations
	mockRepo.On("CreateWalletAccount", mock.Anything, mock.AnythingOfType("*models.WalletAccount")).Return(nil)
	mockClient.On("GetBlockCount", mock.Anything).Return(100, nil)
	mockClient.On("GetBalance", mock.Anything, mock.Anything, mock.Anything).Return("100", nil)
	mockRepo.On("ListTransactions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]models.Transaction{}, 0, nil)
	mockRepo.On("ListWalletsByService", mock.Anything, mock.Anything).Return([]models.WalletAccount{}, nil)
	mockRepo.On("UpdateTransactionStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything, 
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create the transaction service with proper casting
	var clientInterface blockchain.Client = mockClient
	var walletStoreInterface blockchain.WalletStore = mockWalletStore

	// Create pointers to interfaces as required by NewTransactionService
	client := &clientInterface
	store := &walletStoreInterface

	service := blockchain.NewTransactionService(
		mockRepo,
		client,
		store,
		1, // confirmationBlocks
	)

	return mockRepo, mockClient, mockWalletStore, service
}
