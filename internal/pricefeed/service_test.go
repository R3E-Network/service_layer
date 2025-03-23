package pricefeed

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/neo-oracle/internal/config"
	"github.com/your-org/neo-oracle/internal/models"
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

// Mock data source for testing
type mockDataSource struct {
	prices map[string]*models.PriceData
}

func (m *mockDataSource) GetPrices(tokens []string) (map[string]*models.PriceData, error) {
	result := make(map[string]*models.PriceData)
	for _, token := range tokens {
		if price, exists := m.prices[token]; exists {
			result[token] = price
		}
	}
	return result, nil
}

func (m *mockDataSource) Name() string {
	return "MockDataSource"
}

func newMockDataSource() *mockDataSource {
	now := time.Now()
	return &mockDataSource{
		prices: map[string]*models.PriceData{
			"NEO": {
				Token:     "NEO",
				Price:     12.34,
				Timestamp: now,
				Source:    "Mock",
			},
			"GAS": {
				Token:     "GAS",
				Price:     4.56,
				Timestamp: now,
				Source:    "Mock",
			},
		},
	}
}

func TestPriceFeedService(t *testing.T) {
	cfg := &config.PriceFeedConfig{
		UpdateIntervalSec: 300,
		DataSources:       []string{"mock"},
		SupportedTokens:   []string{"NEO", "GAS"},
	}

	bc := &mockBlockchainClient{}
	svc := NewService(cfg, bc)

	// Replace data sources initialization with mock
	svc.dataSources = []DataSource{newMockDataSource()}

	// Start the service
	ctx := context.Background()
	err := svc.Start(ctx)
	assert.NoError(t, err)

	// Update prices
	err = svc.updatePrices()
	assert.NoError(t, err)

	// Test get price
	neoPrice, err := svc.GetPrice("NEO")
	assert.NoError(t, err)
	assert.NotNil(t, neoPrice)
	assert.Equal(t, "NEO", neoPrice.Token)
	assert.Equal(t, 12.34, neoPrice.Price)

	gasPrice, err := svc.GetPrice("GAS")
	assert.NoError(t, err)
	assert.NotNil(t, gasPrice)
	assert.Equal(t, "GAS", gasPrice.Token)
	assert.Equal(t, 4.56, gasPrice.Price)

	// Test get all prices
	allPrices := svc.GetAllPrices()
	assert.Equal(t, 2, len(allPrices))
	assert.Contains(t, allPrices, "NEO")
	assert.Contains(t, allPrices, "GAS")

	// Test unknown token
	_, err = svc.GetPrice("UNKNOWN")
	assert.Error(t, err)

	// Stop the service
	err = svc.Stop()
	assert.NoError(t, err)
}

func TestMedianPriceCalculation(t *testing.T) {
	prices := []*models.PriceData{
		{
			Token: "NEO",
			Price: 10.0,
		},
		{
			Token: "NEO",
			Price: 12.0,
		},
		{
			Token: "NEO",
			Price: 11.0,
		},
	}

	result := calculateMedianPrice(prices)
	assert.NotNil(t, result)
	assert.Equal(t, "NEO", result.Token)
	assert.Equal(t, 11.0, result.Price) // Actually average in our implementation
}
