package service

import (
	"context"
	"log"
	"sync"

	"github.com/R3E-Network/service_layer/internal/blockchain"
	"github.com/R3E-Network/service_layer/internal/common"
	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/pkg/logger"
	"github.com/R3E-Network/service_layer/internal/tee"
)

// ServiceLayer represents the main oracle service coordinator
type ServiceLayer struct {
	config     *config.Config
	blockchain blockchain.BlockchainClient
	tee        *tee.Manager
	services   []common.Service
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// NewServiceLayer creates a new service layer instance
func NewServiceLayer(cfg *config.Config) (*ServiceLayer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize logger
	log := logger.New(logger.LoggingConfig{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		FilePrefix: cfg.Logging.FilePrefix,
	})

	// Initialize blockchain client factory
	bcFactory := blockchain.NewNeoBlockchainClientFactory(&cfg.Neo, log)
	
	// Create blockchain client
	bc, err := bcFactory.Create()
	if err != nil {
		cancel()
		return nil, err
	}

	// Initialize TEE manager
	teeManager, err := tee.NewManager(cfg, log)
	if err != nil {
		cancel()
		return nil, err
	}

	return &ServiceLayer{
		config:     cfg,
		blockchain: bc,
		tee:        teeManager,
		ctx:        ctx,
		cancelFunc: cancel,
		services:   make([]common.Service, 0),
	}, nil
}

// RegisterService adds a service to be managed by the service layer
func (s *ServiceLayer) RegisterService(service common.Service) {
	s.services = append(s.services, service)
}

// Start initiates all registered services
func (s *ServiceLayer) Start() error {
	log.Println("Starting Neo N3 Oracle Service Layer...")

	// Start blockchain client
	if err := s.blockchain.CheckHealth(s.ctx); err != nil {
		return err
	}

	// Initialize TEE environment
	if err := s.tee.Initialize(s.ctx); err != nil {
		return err
	}

	// Start all registered services
	for _, service := range s.services {
		svc := service // Create a copy for the goroutine
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			log.Printf("Starting service: %s", svc.Name())
			if err := svc.Start(s.ctx); err != nil {
				log.Printf("Error starting service %s: %v", svc.Name(), err)
			}
		}()
	}

	return nil
}

// Stop gracefully shuts down all services
func (s *ServiceLayer) Stop() error {
	log.Println("Stopping Neo N3 Oracle Service Layer...")

	// Signal cancellation to all services
	s.cancelFunc()

	// Wait for all services to complete their work
	s.wg.Wait()

	// Stop blockchain client (close connections)
	if err := s.blockchain.Close(); err != nil {
		log.Printf("Error closing blockchain client: %v", err)
	}

	// Cleanup TEE environment
	if err := s.tee.Cleanup(); err != nil {
		log.Printf("Error cleaning up TEE environment: %v", err)
	}

	log.Println("Neo N3 Oracle Service Layer stopped")
	return nil
}

// GetBlockchainClient returns the blockchain client for use by services
func (s *ServiceLayer) GetBlockchainClient() blockchain.BlockchainClient {
	return s.blockchain
}

// GetTEEManager returns the TEE manager for use by services
func (s *ServiceLayer) GetTEEManager() *tee.Manager {
	return s.tee
}
