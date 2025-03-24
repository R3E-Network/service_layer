package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/R3E-Network/service_layer/internal/blockchain/tests/mocks"
	"github.com/R3E-Network/service_layer/internal/models"
)

// BenchmarkTransactionCreation benchmarks the creation of transactions
func BenchmarkTransactionCreation(b *testing.B) {
	// Create mocks
	mockRepo, mockClient, mockWalletStore, service := mocks.CreateMockServices()

	// Setup GetWalletByService
	wallet := &models.WalletAccount{
		ID:                 uuid.New(),
		Service:            "test-service",
		Address:            "NeoAddress123",
		EncryptedPrivateKey: "encrypted-private-key",
		PublicKey:          "public-key",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockRepo.On("GetWalletByService", mock.Anything, "test-service").Return(wallet, nil)

	// Setup CreateTransaction
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

	// Setup for wallet store
	mockWalletStore.On("GetWallet", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, nil)
	
	// Setup for client
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("blockchain.TransactionParams")).Return("tx-hash", nil)
	mockClient.On("GetHeight").Return(uint32(100), nil)

	// Create a common request template
	entityID := uuid.New()
	requestTemplate := models.CreateTransactionRequest{
		Service:    "test-service",
		EntityID:   entityID,
		EntityType: "test-entity",
		Type:       models.TransactionTypeInvoke,
		Script:     "test-script",
		Params:     []interface{}{"param1", "param2"},
		Signers: []models.ScriptSigner{
			{
				Account:          "test-signer",
				Scopes:           "Global",
				AllowedContracts: []string{"contract1", "contract2"},
			},
		},
		Priority:   "high",
	}

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a copy of the request to avoid potential race conditions
		request := requestTemplate
		_, err := service.CreateTransaction(context.Background(), request)
		if err != nil {
			b.Fatalf("Failed to create transaction: %v", err)
		}
	}
}

// TestConcurrentTransactionCreation tests concurrent transaction creation
func TestConcurrentTransactionCreation(t *testing.T) {
	// Create mocks
	mockRepo, mockClient, mockWalletStore, service := mocks.CreateMockServices()

	// Setup GetWalletByService
	wallet := &models.WalletAccount{
		ID:                 uuid.New(),
		Service:            "test-service",
		Address:            "NeoAddress123",
		EncryptedPrivateKey: "encrypted-private-key",
		PublicKey:          "public-key",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockRepo.On("GetWalletByService", mock.Anything, mock.Anything).Return(wallet, nil)

	// Setup CreateTransaction
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

	// Setup for wallet store
	mockWalletStore.On("GetWallet", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, nil)
	
	// Setup for client
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("blockchain.TransactionParams")).Return("tx-hash", nil)
	mockClient.On("GetHeight").Return(uint32(100), nil)

	// Number of concurrent transactions to create
	numConcurrent := 10

	// Create a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	wg.Add(numConcurrent)

	// Create a channel to collect errors
	errorCh := make(chan error, numConcurrent)

	// Create transactions concurrently
	for i := 0; i < numConcurrent; i++ {
		go func(idx int) {
			defer wg.Done()

			// Create a transaction request
			txReq := models.CreateTransactionRequest{
				Service:    "test-service",
				EntityID:   uuid.New(),
				EntityType: "test-entity",
				Type:       models.TransactionTypeInvoke,
				Script:     "test-script",
				Params:     []interface{}{"param1", "param2"},
				Signers: []models.ScriptSigner{
					{
						Account:          "test-signer",
						Scopes:           "Global",
						AllowedContracts: []string{"contract1", "contract2"},
					},
				},
				Priority:   "high",
			}

			// Create the transaction
			_, err := service.CreateTransaction(context.Background(), txReq)
			if err != nil {
				errorCh <- err
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errorCh)

	// Check for errors
	var errors []error
	for err := range errorCh {
		errors = append(errors, err)
	}

	assert.Empty(t, errors, "Expected no errors during concurrent transaction creation")
}

// TestTransactionCreationWithLatency tests transaction creation with simulated network latency
func TestTransactionCreationWithLatency(t *testing.T) {
	// Create mocks
	mockRepo, mockClient, mockWalletStore, service := mocks.CreateMockServices()

	// Setup GetWalletByService
	wallet := &models.WalletAccount{
		ID:                 uuid.New(),
		Service:            "test-service",
		Address:            "NeoAddress123",
		EncryptedPrivateKey: "encrypted-private-key",
		PublicKey:          "public-key",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockRepo.On("GetWalletByService", mock.Anything, mock.Anything).Return(wallet, nil)

	// Setup CreateTransaction
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

	// Setup for wallet store - add simulated latency
	mockWalletStore.On("GetWallet", mock.AnythingOfType("int"), mock.AnythingOfType("string")).
		Run(func(args mock.Arguments) {
			// Simulate network latency
			time.Sleep(50 * time.Millisecond)
		}).
		Return(nil, nil)
	
	// Setup for client - add simulated latency
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("blockchain.TransactionParams")).
		Run(func(args mock.Arguments) {
			// Simulate blockchain network latency
			time.Sleep(100 * time.Millisecond)
		}).
		Return("tx-hash", nil)
		
	mockClient.On("GetHeight").Return(uint32(100), nil)

	// Create a transaction request
	txReq := models.CreateTransactionRequest{
		Service:    "test-service",
		EntityID:   uuid.New(),
		EntityType: "test-entity",
		Type:       models.TransactionTypeInvoke,
		Script:     "test-script",
		Params:     []interface{}{"param1", "param2"},
		Signers: []models.ScriptSigner{
			{
				Account:          "test-signer",
				Scopes:           "Global",
				AllowedContracts: []string{"contract1", "contract2"},
			},
		},
		Priority:   "high",
	}

	// Measure time taken
	start := time.Now()
	txID, err := service.CreateTransaction(context.Background(), txReq)
	elapsed := time.Since(start)

	// Assertions
	require.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, txID)
	
	// The operation should take at least the sum of our simulated latencies
	assert.GreaterOrEqual(t, elapsed, 150*time.Millisecond, "Expected operation to take at least 150ms due to simulated latency")
	
	// Log the time taken for information
	t.Logf("Transaction creation with simulated latency took %v", elapsed)
}
