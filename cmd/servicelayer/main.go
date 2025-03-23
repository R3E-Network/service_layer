package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-org/neo-oracle/internal/config"
	"github.com/your-org/neo-oracle/internal/service"
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

	// Register all services
	if err := serviceLayer.RegisterServices(); err != nil {
		log.Fatalf("Failed to register services: %v", err)
	}

	// Register signal handlers for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Received signal %v, shutting down...", sig)
		if err := serviceLayer.Stop(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	// Start the service
	if err := serviceLayer.Start(); err != nil {
		log.Fatalf("Failed to start service layer: %v", err)
	}

	log.Println("Neo N3 Oracle Service Layer started successfully")

	// Block forever (until signal)
	select {}
}
