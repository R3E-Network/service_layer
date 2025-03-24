package performance

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/R3E-Network/service_layer/internal/database"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/R3E-Network/service_layer/pkg/cache"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// TestDatabasePerformance benchmarks database operations
func BenchmarkDatabasePerformance(b *testing.B) {
	// Skip this test if no database connection is available
	b.Skip("This test requires a real database connection and should be run manually")

	// Setup database connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/service_layer?sslmode=disable")
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Convert sql.DB to sqlx.DB
	dbx := sqlx.NewDb(db, "postgres")
	
	// Setup logger and cache for repository
	zlog := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	cacheManager := &cache.Manager{}  // Initialize as empty - we won't use caching in tests
	
	// Create repositories
	functionRepo := database.NewFunctionRepository(dbx, cacheManager, zlog)
	txRepo := database.NewSQLTransactionRepository(dbx)

	// Test function CRUD operations
	b.Run("FunctionCRUD", func(b *testing.B) {
		benchmarkFunctionCRUD(b, functionRepo)
	})

	// Test transaction CRUD operations
	b.Run("TransactionCRUD", func(b *testing.B) {
		benchmarkTransactionCRUD(b, txRepo)
	})

	// Test complex queries
	b.Run("ComplexQueries", func(b *testing.B) {
		benchmarkComplexQueries(b, functionRepo, txRepo)
	})

	// Test concurrent operations
	b.Run("ConcurrentOperations", func(b *testing.B) {
		benchmarkConcurrentOperations(b, txRepo)
	})
}

// benchmarkFunctionCRUD tests CRUD operations for functions
func benchmarkFunctionCRUD(b *testing.B, repo *database.FunctionRepository) {
	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Generate a unique ID for this iteration
		id := uuid.New()

		// Create a test function
		fn := &database.Function{
			ID:          id.String(),
			Name:        fmt.Sprintf("Test Function %d", i),
			Description: "Test function for benchmarking",
			SourceCode:  "function main(params) { return { result: 'test' }; }",
			UserID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Create the function
		ctx := context.Background()
		err := repo.Create(ctx, fn)
		if err != nil {
			b.Fatalf("Failed to create function: %v", err)
		}

		// Read the function
		readFn, err := repo.GetByID(ctx, fn.ID)
		if err != nil {
			b.Fatalf("Failed to read function: %v", err)
		}

		// Verify
		if readFn.ID != fn.ID {
			b.Fatalf("Function ID mismatch: expected %s, got %s", fn.ID, readFn.ID)
		}

		// Update the function
		readFn.Description = "Updated description"
		err = repo.Update(ctx, readFn)
		if err != nil {
			b.Fatalf("Failed to update function: %v", err)
		}

		// Delete the function
		err = repo.Delete(ctx, fn.ID, 1)
		if err != nil {
			b.Fatalf("Failed to delete function: %v", err)
		}
	}
}

// benchmarkTransactionCRUD tests CRUD operations for transactions
func benchmarkTransactionCRUD(b *testing.B, repo database.TransactionRepository) {
	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Generate a unique ID for this iteration
		id := uuid.New()

		// Create a test transaction
		tx := &models.Transaction{
			ID:          id,
			Hash:        fmt.Sprintf("tx_hash_%d", i),
			Service:     "test_service",
			Status:      models.TransactionStatusPending,
			Type:        "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Create the transaction
		ctx := context.Background()
		err := repo.CreateTransaction(ctx, tx)
		if err != nil {
			b.Fatalf("Failed to create transaction: %v", err)
		}

		// Read the transaction
		readTx, err := repo.GetTransactionByID(ctx, id)
		if err != nil {
			b.Fatalf("Failed to read transaction: %v", err)
		}

		// Verify
		if readTx.ID != id {
			b.Fatalf("Transaction ID mismatch: expected %s, got %s", id, readTx.ID)
		}

		// Update the transaction status
		err = repo.UpdateTransactionStatus(ctx, id, models.TransactionStatusConfirmed, nil, nil, nil, nil, "")
		if err != nil {
			b.Fatalf("Failed to update transaction: %v", err)
		}

		// Delete the transaction
		err = repo.DeleteTransaction(ctx, id)
		if err != nil {
			b.Fatalf("Failed to delete transaction: %v", err)
		}
	}
}

// benchmarkComplexQueries tests complex database queries
func benchmarkComplexQueries(b *testing.B, functionRepo *database.FunctionRepository, transactionRepo database.TransactionRepository) {
	// Create test data for querying
	setupTestData(b, functionRepo, transactionRepo)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Query transactions with filtering and pagination
		ctx := context.Background()
		
		// List transactions with filtering
		service := "test_service"
		status := models.TransactionStatusPending
		page := 1
		limit := 20
		
		// Execute complex query
		_, err := transactionRepo.ListTransactions(ctx, service, status, nil, page, limit)
		if err != nil {
			b.Fatalf("Failed to search transactions: %v", err)
		}
	}
}

// benchmarkConcurrentOperations tests database operations under concurrent load
func benchmarkConcurrentOperations(b *testing.B, repo database.TransactionRepository) {
	// Run the benchmark with concurrency
	b.SetParallelism(10) // 10 concurrent workers
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Generate a unique ID for this iteration
			id := uuid.New()

			// Create a test transaction
			tx := &models.Transaction{
				ID:          id,
				Hash:        fmt.Sprintf("tx_hash_%d", time.Now().UnixNano()),
				Service:     "test_service",
				Status:      models.TransactionStatusPending,
				Type:        "test",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Create
			ctx := context.Background()
			err := repo.CreateTransaction(ctx, tx)
			if err != nil {
				b.Fatalf("Failed to create transaction: %v", err)
			}

			// Read
			_, err = repo.GetTransactionByID(ctx, id)
			if err != nil {
				b.Fatalf("Failed to get transaction: %v", err)
			}

			// Update
			err = repo.UpdateTransactionStatus(ctx, id, models.TransactionStatusConfirmed, nil, nil, nil, nil, "")
			if err != nil {
				b.Fatalf("Failed to update transaction: %v", err)
			}

			// Delete
			err = repo.DeleteTransaction(ctx, id)
			if err != nil {
				b.Fatalf("Failed to delete transaction: %v", err)
			}
		}
	})
}

// setupTestData creates test data for queries
func setupTestData(b *testing.B, functionRepo *database.FunctionRepository, transactionRepo database.TransactionRepository) {
	// Create test transactions
	for i := 0; i < 100; i++ {
		// Create status based on index
		var status models.TransactionStatus
		if i < 30 {
			status = models.TransactionStatusPending
		} else if i < 60 {
			status = models.TransactionStatusConfirmed
		} else {
			status = models.TransactionStatusFailed
		}

		// Create a test transaction
		id := uuid.New()
		tx := &models.Transaction{
			ID:          id,
			Hash:        fmt.Sprintf("tx_hash_%d", i),
			Service:     "test_service",
			Status:      status,
			Type:        "test",
			CreatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt:   time.Now(),
		}

		// Create transaction
		ctx := context.Background()
		err := transactionRepo.CreateTransaction(ctx, tx)
		if err != nil {
			b.Fatalf("Failed to create test transaction: %v", err)
		}
	}

	// Create test functions
	for i := 0; i < 100; i++ {
		// Create a test function
		id := uuid.New()
		idStr := id.String()
		fn := &database.Function{
			ID:          idStr,
			Name:        fmt.Sprintf("Test Function %d", i),
			Description: "Function for query testing",
			SourceCode:  "function main(params) { return { result: 'test' }; }",
			UserID:      1,
			CreatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt:   time.Now(),
		}

		// Create function
		ctx := context.Background()
		err := functionRepo.Create(ctx, fn)
		if err != nil {
			b.Fatalf("Failed to create test function: %v", err)
		}
	}
}