package tests

import (
	"context"
	"testing"
	"time"

	"github.com/R3E-Network/service_layer/internal/blockchain/tests/mocks"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWalletByService(t *testing.T) {
	// Create a mock repository
	mockRepo := new(mocks.MockTransactionRepository)
	
	// Create a wallet for testing
	wallet := &models.WalletAccount{
		ID:                 uuid.New(),
		Service:            "test-service",
		Address:            "test-address",
		EncryptedPrivateKey: "encrypted-private-key",
		PublicKey:          "public-key",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	
	// Set up expectations for the repo
	mockRepo.On("GetWalletByService", mock.Anything, "test-service").Return(wallet, nil)
	
	// Call the method
	gotWallet, err := mockRepo.GetWalletByService(context.Background(), "test-service")
	
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, gotWallet)
	assert.Equal(t, wallet.ID, gotWallet.ID)
	assert.Equal(t, wallet.Service, gotWallet.Service)
	assert.Equal(t, wallet.Address, gotWallet.Address)
	assert.Equal(t, wallet.EncryptedPrivateKey, gotWallet.EncryptedPrivateKey)
	assert.Equal(t, wallet.PublicKey, gotWallet.PublicKey)
	assert.Equal(t, wallet.CreatedAt, gotWallet.CreatedAt)
	assert.Equal(t, wallet.UpdatedAt, gotWallet.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction(t *testing.T) {
	// Setup using the shared mock services
	mockRepo, mockClient, mockWalletStore, service := mocks.CreateMockServices()

	// Create a wallet to use in the test
	wallet := &models.WalletAccount{
		ID:                 uuid.New(),
		Service:            "test-service",
		Address:            "test-address",
		EncryptedPrivateKey: "encrypted-private-key",
		PublicKey:          "public-key",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// For the mock wallet, we can use nil since we're mocking the interface
	var neoWallet interface{} = nil

	// Set up expectations for the repo
	mockRepo.On("GetWalletByService", mock.Anything, "test-service").Return(wallet, nil)
	
	// Set up expectations for walletStore - this matches the WalletStore interface
	mockWalletStore.On("GetWallet", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(neoWallet, nil)

	// Setup client mocks
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("blockchain.TransactionParams")).Return("tx-hash", nil)
	mockClient.On("GetHeight").Return(uint32(100), nil)

	// Setup mock for saving transaction
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

	// Setup transaction creation
	txReq := models.CreateTransactionRequest{
		Service:    "test-service",
		EntityID:   uuid.New(),
		EntityType: "test-entity",
		Type:       models.TransactionTypeInvoke,
		Script:     "test script",
		Params:     []interface{}{"param1", "param2"},
		Signers: []models.ScriptSigner{
			{
				Account:          "test-signer",
				Scopes:           "Global",
				AllowedContracts: []string{"contract1", "contract2"},
			},
		},
		Priority: "high",
	}

	// Call service to create transaction
	txID, err := service.CreateTransaction(context.Background(), txReq)

	// Assertions
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, txID)
	mockRepo.AssertExpectations(t)
	mockClient.AssertExpectations(t)
	mockWalletStore.AssertExpectations(t)
}

func TestGetTransaction(t *testing.T) {
	// Setup using the shared mock services
	mockRepo, _, _, service := mocks.CreateMockServices()

	// Create a transaction for testing
	tx := &models.Transaction{
		ID:         uuid.New(),
		Service:    "test-service",
		Status:     models.TransactionStatusPending,
		EntityID:   uuid.New(),
		EntityType: "test-entity",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Setup expectations
	mockRepo.On("GetTransactionByID", mock.Anything, tx.ID).Return(tx, nil)

	// Call service
	result, err := service.GetTransaction(context.Background(), tx.ID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, tx.ID, result.ID)
	assert.Equal(t, tx.Status, result.Status)
	mockRepo.AssertExpectations(t)
}

func TestListTransactions(t *testing.T) {
	// Setup using the shared mock services
	mockRepo, _, _, service := mocks.CreateMockServices()

	// Setup test data
	transactions := []models.Transaction{
		{
			ID:         uuid.New(),
			Service:    "test-service",
			Status:     models.TransactionStatusPending,
			EntityID:   uuid.New(),
			EntityType: "test-entity",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			Service:    "test-service",
			Status:     models.TransactionStatusPending,
			EntityID:   uuid.New(),
			EntityType: "test-entity",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	// Setup expectations
	mockRepo.On("ListTransactions", 
		mock.Anything, 
		"test-service", 
		models.TransactionStatusPending, 
		mock.AnythingOfType("*uuid.UUID"), 
		int64(1), 
		int64(10),
	).Return(transactions, int64(len(transactions)), nil)

	// Call service
	var entityID *uuid.UUID = nil
	result, count, err := service.ListTransactions(context.Background(), "test-service", models.TransactionStatusPending, entityID, 1, 10)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	assert.Len(t, result, 2)
	assert.Equal(t, transactions[0].ID, result[0].ID)
	assert.Equal(t, transactions[1].ID, result[1].ID)
	mockRepo.AssertExpectations(t)
}
