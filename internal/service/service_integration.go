package service

import (
	"github.com/your-org/neo-oracle/internal/api"
	"github.com/your-org/neo-oracle/internal/functions"
	"github.com/your-org/neo-oracle/internal/gasbank"
	"github.com/your-org/neo-oracle/internal/logging"
	"github.com/your-org/neo-oracle/internal/metrics"
	"github.com/your-org/neo-oracle/internal/pricefeed"
	"github.com/your-org/neo-oracle/internal/secrets"
	"github.com/your-org/neo-oracle/internal/triggers"
)

// RegisterServices sets up all service components
func (s *ServiceLayer) RegisterServices() error {
	// Initialize logging first
	loggingSvc, err := logging.NewService(&s.config.Logging)
	if err != nil {
		return err
	}
	s.RegisterService(loggingSvc)

	// Initialize metrics
	metricsSvc := metrics.NewService(&s.config.Metrics)
	s.RegisterService(metricsSvc)

	// Initialize secrets service
	secretsSvc := secrets.NewService(s.config.TEE)
	s.RegisterService(secretsSvc)

	// Initialize function executor with JS runtime support
	jsRuntime := functions.NewJSRuntime(secretsSvc, 30) // 30 seconds timeout
	executor := functions.NewExecutor(s.tee)

	// Initialize services
	gasBankSvc := gasbank.NewService(&s.config.GasBank, s.blockchain)
	priceFeedSvc := pricefeed.NewService(&s.config.PriceFeed, s.blockchain)
	triggersSvc := triggers.NewService(executor, priceFeedSvc)

	// Create API server
	apiServer := api.NewServer(&s.config.Server, gasBankSvc, priceFeedSvc, s.tee)

	// Register all services
	s.RegisterService(secretsSvc)
	s.RegisterService(gasBankSvc)
	s.RegisterService(priceFeedSvc)
	s.RegisterService(triggersSvc)
	s.RegisterService(apiServer)

	return nil
}
