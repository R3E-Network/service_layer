package pricefeed

import (
	"errors"
	"testing"
	"time"

	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/R3E-Network/service_layer/internal/blockchain"
	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/pkg/logger"
	"github.com/stretchr/testify/require"
)

// Mock PriceFeedRepository
type MockPriceFeedRepository struct {
	mock.Mock
}

func (m *MockPriceFeedRepository) CreatePriceFeed(ctx interface{}, feed *models.PriceFeed) (*models.PriceFeed, error) {
	args := m.Called(ctx, feed)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) GetPriceFeedByID(ctx interface{}, id string) (*models.PriceFeed, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) GetPriceFeedByPair(ctx interface{}, pair string) (*models.PriceFeed, error) {
	args := m.Called(ctx, pair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) GetPriceFeed(ctx interface{}, id string) (*models.PriceFeed, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) ListPriceFeeds(ctx interface{}) ([]*models.PriceFeed, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) UpdatePriceFeed(ctx interface{}, feed *models.PriceFeed) (*models.PriceFeed, error) {
	args := m.Called(ctx, feed)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) DeletePriceFeed(ctx interface{}, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPriceFeedRepository) CreatePriceData(ctx interface{}, data *models.PriceData) (*models.PriceData, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceData), args.Error(1)
}

func (m *MockPriceFeedRepository) GetLatestPriceData(ctx interface{}, priceFeedID string) (*models.PriceData, error) {
	args := m.Called(ctx, priceFeedID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceData), args.Error(1)
}

func (m *MockPriceFeedRepository) GetPriceDataHistory(ctx interface{}, priceFeedID string, limit int, offset int) ([]*models.PriceData, error) {
	args := m.Called(ctx, priceFeedID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PriceData), args.Error(1)
}

func (m *MockPriceFeedRepository) CreatePriceSource(ctx interface{}, source *models.PriceSource) (*models.PriceSource, error) {
	args := m.Called(ctx, source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceSource), args.Error(1)
}

func (m *MockPriceFeedRepository) GetPriceSourceByID(ctx interface{}, id int) (*models.PriceSource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceSource), args.Error(1)
}

func (m *MockPriceFeedRepository) GetPriceSourceByName(ctx interface{}, name string) (*models.PriceSource, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceSource), args.Error(1)
}

func (m *MockPriceFeedRepository) ListPriceSources(ctx interface{}) ([]*models.PriceSource, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PriceSource), args.Error(1)
}

func (m *MockPriceFeedRepository) UpdatePriceSource(ctx interface{}, source *models.PriceSource) (*models.PriceSource, error) {
	args := m.Called(ctx, source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceSource), args.Error(1)
}

func (m *MockPriceFeedRepository) DeletePriceSource(ctx interface{}, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPriceFeedRepository) GetPriceFeedBySymbol(ctx interface{}, symbol string) (*models.PriceFeed, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceFeed), args.Error(1)
}

func (m *MockPriceFeedRepository) AddPriceSource(ctx interface{}, priceFeedID string, source *models.PriceSource) (*models.PriceSource, error) {
	args := m.Called(ctx, priceFeedID, source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceSource), args.Error(1)
}

func (m *MockPriceFeedRepository) RemovePriceSource(ctx interface{}, priceFeedID string, sourceID string) error {
	args := m.Called(ctx, priceFeedID, sourceID)
	return args.Error(0)
}

func (m *MockPriceFeedRepository) UpdateLastPrice(ctx interface{}, id string, price float64, txHash string) error {
	args := m.Called(ctx, id, price, txHash)
	return args.Error(0)
}

func (m *MockPriceFeedRepository) RecordPriceUpdate(ctx interface{}, id string, price float64, txID string) error {
	args := m.Called(ctx, id, price, txID)
	return args.Error(0)
}

func (m *MockPriceFeedRepository) GetPriceHistory(ctx interface{}, id string, limit int) ([]*models.PriceUpdate, error) {
	args := m.Called(ctx, id, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PriceUpdate), args.Error(1)
}

// Mock PriceFetcherFactory
type MockPriceFetcherFactory struct {
	mock.Mock
}

func (m *MockPriceFetcherFactory) CreateFetcher(source *models.PriceSource) (PriceFetcher, error) {
	args := m.Called(source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(PriceFetcher), args.Error(1)
}

// Mock PriceFetcher
type MockPriceFetcher struct {
	mock.Mock
}

func (m *MockPriceFetcher) FetchPrice(ctx interface{}, baseToken, quoteToken string) (float64, error) {
	args := m.Called(ctx, baseToken, quoteToken)
	return args.Get(0).(float64), args.Error(1)
}

// Mock PriceAggregator
type MockPriceAggregator struct {
	mock.Mock
}

func (m *MockPriceAggregator) Aggregate(prices map[string]float64, weights map[string]float64) (float64, error) {
	args := m.Called(prices, weights)
	return args.Get(0).(float64), args.Error(1)
}

// Mock GasBankService
type MockGasBankService struct {
	mock.Mock
}

func (m *MockGasBankService) AllocateGas(ctx interface{}, userID int, operation string, estimatedGas int64) (int64, error) {
	args := m.Called(ctx, userID, operation, estimatedGas)
	return args.Get(0).(int64), args.Error(1)
}

// Mock TEEManager
type MockTEEManager struct {
	mock.Mock
}

// Helper function to setup test service
func setupTestService() (*PriceFeedService, *MockPriceFeedRepository, *MockPriceFetcherFactory, *MockPriceAggregator, *MockGasBankService, *MockTEEManager) {
	cfg := &config.Config{
		Services: config.ServicesConfig{
			PriceFeed: config.PriceFeedApi{
				NumWorkers: 1,
			},
		},
	}
	log := logger.NewLogger("test")
	mockRepo := new(MockPriceFeedRepository)
	mockBlockchainClient := &blockchain.Client{}
	mockGasBankService := new(MockGasBankService)
	mockTEEManager := new(MockTEEManager)

	service := NewService(cfg, log, mockRepo, mockBlockchainClient, mockGasBankService, mockTEEManager)

	mockFetcherFactory := new(MockPriceFetcherFactory)
	mockAggregator := new(MockPriceAggregator)
	
	// Set mock implementations directly in the service
	service.fetcherFactory = mockFetcherFactory
	service.aggregator = mockAggregator

	return service, mockRepo, mockFetcherFactory, mockAggregator, mockGasBankService, mockTEEManager
}

// TestCreatePriceFeed tests price feed creation
func TestCreatePriceFeed(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceFeedByPair", mock.Anything, "BTC/USD").Return(nil, errors.New("not found")).Once()

		expectedFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		mockRepo.On("CreatePriceFeed", mock.Anything, mock.AnythingOfType("*models.PriceFeed")).Return(expectedFeed, nil).Once()

		// Call service
		feed, err := service.CreatePriceFeed("BTC", "USD", "1h", 0.5, "24h", "0x1234")

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, expectedFeed, feed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("PairAlreadyExists", func(t *testing.T) {
		// Setup mock
		existingFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
		}
		mockRepo.On("GetPriceFeedByPair", mock.Anything, "BTC/USD").Return(existingFeed, nil).Once()

		// Call service
		feed, err := service.CreatePriceFeed("BTC", "USD", "1h", 0.5, "24h", "0x1234")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		// Call service with invalid input
		feed, err := service.CreatePriceFeed("", "USD", "1h", 0.5, "24h", "0x1234")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		assert.Contains(t, err.Error(), "required")

		// No repository calls should be made
		mockRepo.AssertNotCalled(t, "GetPriceFeedByPair")
		mockRepo.AssertNotCalled(t, "CreatePriceFeed")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceFeedByPair", mock.Anything, "BTC/USD").Return(nil, errors.New("not found")).Once()
		mockRepo.On("CreatePriceFeed", mock.Anything, mock.AnythingOfType("*models.PriceFeed")).Return(nil, errors.New("database error")).Once()

		// Call service
		feed, err := service.CreatePriceFeed("BTC", "USD", "1h", 0.5, "24h", "0x1234")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdatePriceFeed tests price feed update
func TestUpdatePriceFeed(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		existingFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
		}
		mockRepo.On("GetPriceFeedByID", mock.Anything, "1").Return(existingFeed, nil).Once()

		updatedFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "2h",
			DeviationThreshold: 0.8,
			HeartbeatInterval:  "48h",
			ContractAddress:    "0x5678",
			Active:             false,
		}
		mockRepo.On("UpdatePriceFeed", mock.Anything, mock.AnythingOfType("*models.PriceFeed")).Return(updatedFeed, nil).Once()

		// Call service
		feed, err := service.UpdatePriceFeed("1", "BTC", "USD", "2h", 0.8, "48h", "0x5678", false)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, updatedFeed, feed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FeedNotFound", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceFeedByID", mock.Anything, "999").Return(nil, errors.New("not found")).Once()

		// Call service
		feed, err := service.UpdatePriceFeed("999", "BTC", "USD", "2h", 0.8, "48h", "0x5678", false)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		// Call service with invalid input
		feed, err := service.UpdatePriceFeed("1", "", "USD", "2h", 0.8, "48h", "0x5678", false)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		assert.Contains(t, err.Error(), "required")

		// No repository calls should be made
		mockRepo.AssertNotCalled(t, "GetPriceFeedByID")
		mockRepo.AssertNotCalled(t, "UpdatePriceFeed")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Setup mock
		existingFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
		}
		mockRepo.On("GetPriceFeedByID", mock.Anything, "1").Return(existingFeed, nil).Once()
		mockRepo.On("UpdatePriceFeed", mock.Anything, mock.AnythingOfType("*models.PriceFeed")).Return(nil, errors.New("database error")).Once()

		// Call service
		feed, err := service.UpdatePriceFeed("1", "BTC", "USD", "2h", 0.8, "48h", "0x5678", false)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

// TestDeletePriceFeed tests price feed deletion
func TestDeletePriceFeed(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		mockRepo.On("DeletePriceFeed", mock.Anything, "1").Return(nil).Once()

		// Call service
		err := service.DeletePriceFeed("1")

		// Assertions
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Setup mock
		mockRepo.On("DeletePriceFeed", mock.Anything, "1").Return(errors.New("database error")).Once()

		// Call service
		err := service.DeletePriceFeed("1")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

// TestGetPriceFeed tests retrieving a price feed
func TestGetPriceFeed(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		expectedFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
		}
		mockRepo.On("GetPriceFeedByID", mock.Anything, "1").Return(expectedFeed, nil).Once()

		// Call service
		feed, err := service.GetPriceFeed("1")

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, expectedFeed, feed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceFeedByID", mock.Anything, "999").Return(nil, errors.New("not found")).Once()

		// Call service
		feed, err := service.GetPriceFeed("999")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetPriceFeedByPair tests retrieving a price feed by pair
func TestGetPriceFeedByPair(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		expectedFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
		}
		mockRepo.On("GetPriceFeedByPair", mock.Anything, "BTC/USD").Return(expectedFeed, nil).Once()

		// Call service
		feed, err := service.GetPriceFeedByPair("BTC/USD")

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, expectedFeed, feed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceFeedByPair", mock.Anything, "XYZ/ABC").Return(nil, errors.New("not found")).Once()

		// Call service
		feed, err := service.GetPriceFeedByPair("XYZ/ABC")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feed)
		mockRepo.AssertExpectations(t)
	})
}

// TestListPriceFeeds tests listing all price feeds
func TestListPriceFeeds(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		expectedFeeds := []*models.PriceFeed{
			{
				ID:                 "1",
				BaseToken:          "BTC",
				QuoteToken:         "USD",
				Pair:               "BTC/USD",
				UpdateInterval:     "1h",
				DeviationThreshold: 0.5,
				HeartbeatInterval:  "24h",
				ContractAddress:    "0x1234",
				Active:             true,
			},
			{
				ID:                 "2",
				BaseToken:          "ETH",
				QuoteToken:         "USD",
				Pair:               "ETH/USD",
				UpdateInterval:     "1h",
				DeviationThreshold: 0.5,
				HeartbeatInterval:  "24h",
				ContractAddress:    "0x5678",
				Active:             true,
			},
		}
		mockRepo.On("ListPriceFeeds", mock.Anything).Return(expectedFeeds, nil).Once()

		// Call service
		feeds, err := service.ListPriceFeeds()

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, expectedFeeds, feeds)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Setup mock
		mockRepo.On("ListPriceFeeds", mock.Anything).Return(nil, errors.New("database error")).Once()

		// Call service
		feeds, err := service.ListPriceFeeds()

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, feeds)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetLatestPrice tests retrieving the latest price
func TestGetLatestPrice(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		expectedPrice := &models.PriceData{
			ID:          "1",
			PriceFeedID: "1",
			Price:       50000.0,
			Timestamp:   time.Now(),
			RoundID:     "123",
			TxHash:      "0xabcd",
			Source:      "test",
		}
		mockRepo.On("GetLatestPriceData", mock.Anything, "1").Return(expectedPrice, nil).Once()

		// Call service
		price, err := service.GetLatestPrice("1")

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, expectedPrice, price)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetLatestPriceData", mock.Anything, "999").Return(nil, errors.New("not found")).Once()

		// Call service
		price, err := service.GetLatestPrice("999")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, price)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetPriceHistory tests retrieving price history
func TestGetPriceHistory(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		expectedPrices := []*models.PriceData{
			{
				ID:          "1",
				PriceFeedID: "1",
				Price:       50000.0,
				Timestamp:   time.Now().Add(-time.Hour),
				RoundID:     "122",
				TxHash:      "0xabcd1",
				Source:      "test",
			},
			{
				ID:          "2",
				PriceFeedID: "1",
				Price:       51000.0,
				Timestamp:   time.Now(),
				RoundID:     "123",
				TxHash:      "0xabcd2",
				Source:      "test",
			},
		}
		mockRepo.On("GetPriceDataHistory", mock.Anything, "1", 10, 0).Return(expectedPrices, nil).Once()

		// Call service
		prices, err := service.GetPriceHistory("1", 10, 0)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, expectedPrices, prices)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceDataHistory", mock.Anything, "1", 10, 0).Return(nil, errors.New("database error")).Once()

		// Call service
		prices, err := service.GetPriceHistory("1", 10, 0)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, prices)
		mockRepo.AssertExpectations(t)
	})
}

// TestTriggerPriceUpdate tests manually triggering a price update
func TestTriggerPriceUpdate(t *testing.T) {
	service, mockRepo, _, _, _, _ := setupTestService()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		existingFeed := &models.PriceFeed{
			ID:                 "1",
			BaseToken:          "BTC",
			QuoteToken:         "USD",
			Pair:               "BTC/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x1234",
			Active:             true,
		}
		mockRepo.On("GetPriceFeedByID", mock.Anything, "1").Return(existingFeed, nil).Once()

		// Call service
		err := service.TriggerPriceUpdate("1")

		// Assertions
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FeedNotFound", func(t *testing.T) {
		// Setup mock
		mockRepo.On("GetPriceFeedByID", mock.Anything, "999").Return(nil, errors.New("not found")).Once()

		// Call service
		err := service.TriggerPriceUpdate("999")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InactiveFeed", func(t *testing.T) {
		// Setup mock
		inactiveFeed := &models.PriceFeed{
			ID:                 "2",
			BaseToken:          "ETH",
			QuoteToken:         "USD",
			Pair:               "ETH/USD",
			UpdateInterval:     "1h",
			DeviationThreshold: 0.5,
			HeartbeatInterval:  "24h",
			ContractAddress:    "0x5678",
			Active:             false,
		}
		mockRepo.On("GetPriceFeedByID", mock.Anything, "2").Return(inactiveFeed, nil).Once()

		// Call service
		err := service.TriggerPriceUpdate("2")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "inactive")
		mockRepo.AssertExpectations(t)
	})
}
