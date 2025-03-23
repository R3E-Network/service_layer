package service

import (
	"context"
	"log"
	"sync"

	"github.com/your-org/neo-oracle/internal/blockchain"
	"github.com/your-org/neo-oracle/internal/common"
	"github.com/your-org/neo-oracle/internal/config"
	"github.com/your-org/neo-oracle/internal/tee"
)

// ServiceLayer represents the main oracle service coordinator
type ServiceLayer struct {
	config     *config.Config
	blockchain *blockchain.Client
	tee        *tee.Manager
	services   []common.Service
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// NewServiceLayer creates a new service layer instance
func NewServiceLayer(cfg *config.Config) (*ServiceLayer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize blockchain client
	bc, err := blockchain.NewClient(cfg.Blockchain)
	if err != nil {
		cancel()
		return nil, err
	}

	// Initialize TEE manager
	teeManager, err := tee.NewManager(cfg.TEE)
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
	if err := s.blockchain.Connect(); err != nil {
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
	s.cancelFunc()

	// Wait for all services to stop
	s.wg.Wait()

	// Stop TEE manager
	if err := s.tee.Shutdown(); err != nil {
		log.Printf("Error shutting down TEE: %v", err)
	}

	// Disconnect blockchain client
	if err := s.blockchain.Disconnect(); err != nil {
		log.Printf("Error disconnecting blockchain client: %v", err)
	}

	return nil
}
