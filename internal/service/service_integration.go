package service

import (
	"github.com/R3E-Network/service_layer/internal/api"
	"github.com/R3E-Network/service_layer/internal/functions"
	"github.com/R3E-Network/service_layer/internal/logging"
	"github.com/R3E-Network/service_layer/internal/metrics"
	"github.com/R3E-Network/service_layer/internal/secrets"
	"github.com/R3E-Network/service_layer/internal/triggers"
	"github.com/R3E-Network/service_layer/pkg/logger"
	// Commenting out unused imports for now
	// "github.com/R3E-Network/service_layer/internal/gasbank"
	// "github.com/R3E-Network/service_layer/internal/pricefeed"
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
	secretsSvc := secrets.NewService(&s.config.TEE)
	s.RegisterService(secretsSvc)

	// Initialize function executor with JS runtime support
	jsRuntime := functions.NewJSRuntime(secretsSvc, 30) // 30 seconds timeout
	executor := functions.NewExecutor(s.tee)

	// For now, we'll comment out the more complex service initializations that require
	// repositories until we can properly implement them
	/*
	// Initialize services
	gasBankSvc, err := gasbank.NewService(s.config, nil, s.blockchain, s.tee)
	if err != nil {
		return err
	}
	
	priceFeedSvc, err := pricefeed.NewService(s.config, nil, s.blockchain, s.tee)
	if err != nil {
		return err
	}
	*/
	
	// Create simple trigger service
	triggersSvc := triggers.NewService(executor, nil)

	// Create API server - create a new logger
	apiLogger := logger.New(logger.LoggingConfig{
		Level:      s.config.Logging.Level,
		Format:     s.config.Logging.Format,
		Output:     s.config.Logging.Output,
		FilePrefix: s.config.Logging.FilePrefix,
	})
	
	// Note: We're not using the apiServer for now since it doesn't implement the Service interface
	// Ignoring the return value with an underscore to prevent unused variable error
	_, err = api.NewServer(s.config, apiLogger)
	if err != nil {
		return err
	}

	// Register all services
	s.RegisterService(secretsSvc)
	// s.RegisterService(gasBankSvc)
	// s.RegisterService(priceFeedSvc)
	s.RegisterService(triggersSvc)
	
	// For now, we'll skip registering the API server since it doesn't implement
	// the Service interface (missing Name method)
	// We'll need to implement this properly in the future
	// s.RegisterService(apiServer)

	return nil
}
