# Trusted Execution Environment (TEE) Documentation

## Overview

The Trusted Execution Environment (TEE) module provides secure execution capabilities for running JavaScript functions in isolated, attested environments. It ensures that sensitive code and data are protected from unauthorized access, even from privileged users or system administrators.

## Architecture

The TEE implementation follows a provider-based architecture:

1. **Manager**: Coordinates TEE operations and interfaces with the rest of the system
2. **Provider Interface**: Defines the contract that all TEE providers must implement
3. **Providers**: Concrete implementations for different TEE platforms (currently Azure Confidential Computing)
4. **Runtime**: Provides the JavaScript execution environment within the TEE

### Component Relationships

```
┌────────────────┐     ┌────────────────┐     ┌────────────────┐
│                │     │                │     │                │
│    Service     │────▶│    Manager     │────▶│    Provider    │
│                │     │                │     │    Interface   │
└────────────────┘     └────────────────┘     └────────────────┘
                                                      │
                                                      │
                                                      ▼
                                            ┌────────────────────┐
                                            │                    │
                                            │ Azure Provider     │
                                            │                    │
                                            └────────────────────┘
```

## Configuration

The TEE functionality is configured through the application's configuration system:

```go
type TEEConfig struct {
    Enabled             bool          `mapstructure:"enabled"`
    Provider            string        `mapstructure:"provider"`
    AttestationEndpoint string        `mapstructure:"attestation_endpoint"`
    Azure               AzureConfig   `mapstructure:"azure"`
    Runtime             RuntimeConfig `mapstructure:"runtime"`
}

type AzureConfig struct {
    ClientID                string `mapstructure:"client_id"`
    ClientSecret            string `mapstructure:"client_secret"`
    TenantID                string `mapstructure:"tenant_id"`
    AttestationEndpoint     string `mapstructure:"attestation_endpoint"`
    AttestationPolicyName   string `mapstructure:"attestation_policy_name"`
    ConfidentialLedgerURL   string `mapstructure:"confidential_ledger_url"`
    KeyVaultURL             string `mapstructure:"key_vault_url"`
    Runtime                 RuntimeConfig `mapstructure:"runtime"`
}

type RuntimeConfig struct {
    JSMemoryLimit      int `mapstructure:"js_memory_limit"`
    ExecutionTimeout   int `mapstructure:"execution_timeout"`
}
```

### Configuration Options

- **enabled**: Enables/disables TEE functionality
- **provider**: The TEE provider to use ("azure" currently supported)
- **attestation_endpoint**: Endpoint for attestation service
- **azure**: Azure-specific configuration
  - **client_id**: Azure AD application client ID
  - **client_secret**: Azure AD application client secret
  - **tenant_id**: Azure AD tenant ID
  - **attestation_endpoint**: Azure attestation service endpoint
  - **attestation_policy_name**: Custom attestation policy name
  - **confidential_ledger_url**: URL for Azure Confidential Ledger
  - **key_vault_url**: URL for Azure Key Vault
- **runtime**: Runtime configuration
  - **js_memory_limit**: Memory limit for JavaScript runtime in MB
  - **execution_timeout**: Execution timeout in seconds

## Usage

### Function Execution

To execute a function in a TEE:

```go
// Get the TEE manager
teeManager, err := tee.NewManager(config, logger)
if err != nil {
    // Handle error
}

// Execute a function securely
result, err := teeManager.ExecuteSecureFunction(ctx, function, params, secrets)
if err != nil {
    // Handle error
}
```

### Secret Management

Secrets can be securely stored and retrieved in the TEE:

```go
// Store a secret
secret := &models.Secret{
    UserID: userID,
    Name:   "api_key",
    Value:  "secret_value",
}
err := teeManager.StoreSecret(ctx, secret)

// Retrieve a secret
value, err := teeManager.GetSecret(ctx, userID, "api_key")

// Delete a secret
err := teeManager.DeleteSecret(ctx, userID, "api_key")
```

### Attestation

Attestation proves that code is running in a genuine TEE:

```go
// Get attestation report
attestation, err := teeManager.GetAttestation(ctx)
```

## Azure Confidential Computing Provider

The Azure provider implements TEE functionality using Azure Confidential Computing, which leverages Intel SGX enclaves for hardware-level isolation. It provides:

1. **Hardware isolation**: Execute code in protected memory regions
2. **Memory encryption**: All data in memory is encrypted
3. **Attestation**: Verify the environment is genuine and unmodified
4. **Secret protection**: Encrypt and seal secrets to the enclave

## Function Model Extension

The Function model includes fields to support TEE execution:

```go
type Function struct {
    // Existing fields
    ID             int
    UserID         int
    Name           string
    Description    string
    SourceCode     string
    // TEE-specific fields
    Code           string  // Used by TEE for execution
    SecureExecution bool    // Whether to use TEE for execution
    // Other fields
}
```

## Security Considerations

1. **Attestation validation**: Always verify attestation reports before trusting a TEE
2. **Secret management**: Use the TEE's secure storage capabilities for secrets
3. **Code integrity**: Ensure that code hasn't been tampered with
4. **Limited attack surface**: Minimize the code running in the TEE to reduce vulnerabilities
5. **Memory limitations**: Be aware of the memory constraints in TEE environments

## Testing

The TEE implementation includes simulation capabilities for testing:

1. When TEE is disabled, functions run in simulation mode
2. Mock providers can be implemented for testing
3. Tests should verify both functional correctness and security properties

## Error Handling

Common errors include:

- `ErrTEENotEnabled`: TEE functionality is disabled
- `ErrInvalidAttestationReport`: Attestation report is invalid
- `ErrSecureExecutionFailed`: Secure execution failed
- `ErrMemoryLimitExceeded`: Memory limit exceeded during execution
- `ErrExecutionTimeout`: Execution timed out
- `ErrProviderNotInitialized`: Provider not initialized
- `ErrUnsupportedProvider`: Unsupported TEE provider
