package tests

import (
	"context"
	"os"
	"testing"

	"github.com/R3E-Network/service_layer/internal/blockchain"
	"github.com/R3E-Network/service_layer/internal/blockchain/tests/mocks"
	"github.com/R3E-Network/service_layer/internal/database"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestIntegrationCreateTransaction performs an integration test of the transaction management system
// Note: This requires a PostgreSQL database to be running
// To run this test: TEST_DB_DSN="postgres://user:password@localhost:5432/testdb?sslmode=disable" go test -tags=integration
func TestIntegrationCreateTransaction(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get database connection
	dsn := getDatabaseDSN()
	if dsn == "" {
		t.Skip("Skipping integration test, no database DSN provided")
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Initialize database schema (if needed for test)
	// For a real test, you might want to create necessary tables
	require.NoError(t, err)

	// Create repository with real database
	repo := database.NewSQLTransactionRepository(db)

	// Create mocks for client and wallet store
	mockRepo, mockClient, mockWalletStore, _ := mocks.CreateMockServices()

	// Setup test data
	walletAccount := &models.WalletAccount{
		ID:                 uuid.New(),
		Service:            "test-service",
		Address:            "test-address",
		EncryptedPrivateKey: "encrypted-key",
		PublicKey:          "public-key",
	}

	// Mocking repository behavior for wallet retrieval
	mockRepo.On("GetWalletByService", mock.Anything, "test-service").Return(walletAccount, nil)

	// Mocking client behavior for transaction creation
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("blockchain.TransactionParams")).Return("tx-hash", nil)
	mockClient.On("GetHeight").Return(uint32(100), nil)

	// Mock the wallet
	mockWallet := &wallet.Wallet{}
	mockWalletStore.On("GetWallet", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockWallet, nil)

	// Mock transaction repository for CreateTransaction
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

	// Create transaction params - use the correct structure as expected by the service
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
		Priority: "high",
	}

	// Use the real repo for this integration test
	// Cast and create a service with the real repo
	var client blockchain.Client = mockClient
	var store blockchain.WalletStore = mockWalletStore
	service := blockchain.NewTransactionService(
		repo,
		&client,
		&store,
		1, // confirmationBlocks
	)

	// Call the service with the correct request type
	txID, err := service.CreateTransaction(context.Background(), txReq)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, txID)
	assert.NotEqual(t, uuid.Nil, txID)

	// No need to verify mock expectations for repo since we're using a real database
	// But we do verify the client and wallet store mocks
	mockClient.AssertExpectations(t)
	mockWalletStore.AssertExpectations(t)
}

// Helper function to get database DSN from environment
func getDatabaseDSN() string {
	return os.Getenv("TEST_DB_DSN")
}
