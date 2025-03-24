# Service Layer

The Service Layer is the core component that manages and coordinates all services within the application. It provides a unified interface for registering, starting, and stopping services.

## Architecture

The Service Layer follows a modular design where individual services implement the `common.Service` interface and are registered with the service layer. The service layer then manages the lifecycle of these services, including:

1. Starting services in the correct order
2. Stopping services gracefully on shutdown
3. Managing shared resources between services
4. Monitoring service health

## Service Registration

Services are registered with the service layer in the main application entry point:

```go
// Create service layer
serviceLayer, err := service.NewServiceLayer(cfg)
if err != nil {
    log.Fatalf("Failed to create service layer: %v", err)
}

// Register services
priceFeedService, err := pricefeed.NewService(cfg, priceFeedRepo, bcClient, teeManager)
if err != nil {
    log.Fatalf("Failed to create price feed service: %v", err)
}
priceFeedServiceImpl := pricefeed.NewServiceImpl(priceFeedService, cfg)
serviceLayer.RegisterService(priceFeedServiceImpl)

// Similar for other services...
```

## Trusted Execution Environment (TEE) Integration

The Service Layer integrates with the Trusted Execution Environment (TEE) to provide secure execution of sensitive operations. The TEE Manager is initialized during service layer creation and is made available to services that require secure execution.

Key TEE features:
- Support for multiple TEE providers (Azure ACC, AWS Nitro, etc.)
- Environment verification and attestation
- Secure key storage and management
- Isolated execution of sensitive code

## Blockchain Integration

The Service Layer initializes a blockchain client that services can use to interact with the blockchain. This client provides:
- Transaction submission and monitoring
- Smart contract interaction
- Event monitoring
- Wallet management

## Configuration

The Service Layer is configured through the main application configuration, which includes settings for:
- Server configuration
- Blockchain settings
- TEE configuration
- Individual service configurations
