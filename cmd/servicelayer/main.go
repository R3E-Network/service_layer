package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/R3E-Network/service_layer/internal/blockchain"
	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/internal/gasbank"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/R3E-Network/service_layer/internal/pricefeed"
	"github.com/R3E-Network/service_layer/internal/service"
	"github.com/R3E-Network/service_layer/internal/tee"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	createDefaultConfig := flag.Bool("create-config", false, "Create default configuration file")
	flag.Parse()

	// Create default configuration if requested
	if *createDefaultConfig {
		cfg := config.DefaultConfig()
		if err := config.SaveConfig(cfg, *configPath); err != nil {
			log.Fatalf("Failed to create default configuration: %v", err)
		}
		log.Printf("Default configuration created at %s", *configPath)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create service layer
	serviceLayer, err := service.NewServiceLayer(cfg)
	if err != nil {
		log.Fatalf("Failed to create service layer: %v", err)
	}

	// Get references to shared components
	bcClient := serviceLayer.GetBlockchainClient()
	teeManager := serviceLayer.GetTEEManager()

	// Initialize repository instances (in a real implementation, these would connect to a database)
	priceFeedRepo := &models.MockPriceFeedRepository{}
	gasBankRepo := &models.MockGasBankRepository{}

	// Create and register services
	
	// PriceFeed Service
	priceFeedService, err := pricefeed.NewService(cfg, priceFeedRepo, bcClient, teeManager)
	if err != nil {
		log.Fatalf("Failed to create price feed service: %v", err)
	}
	priceFeedServiceImpl := pricefeed.NewServiceImpl(priceFeedService, cfg)
	serviceLayer.RegisterService(priceFeedServiceImpl)
	
	// GasBank Service
	gasBankService, err := gasbank.NewService(cfg, gasBankRepo, bcClient, teeManager)
	if err != nil {
		log.Fatalf("Failed to create gas bank service: %v", err)
	}
	gasBankServiceImpl := gasbank.NewServiceImpl(gasBankService, cfg)
	serviceLayer.RegisterService(gasBankServiceImpl)

	// Register signal handlers for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Received signal %v, shutting down...", sig)
		if err := serviceLayer.Stop(); err != nil {
			log.Printf("Error stopping service layer: %v", err)
		}
		os.Exit(0)
	}()

	// Start the service layer
	if err := serviceLayer.Start(); err != nil {
		log.Fatalf("Failed to start service layer: %v", err)
	}

	// Block forever (the signal handler will exit the program)
	select {}
}
