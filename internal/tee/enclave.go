package tee

import (
	"sync"
	"time"
)

// EnclaveStatus represents the status of an enclave
type EnclaveStatus string

const (
	// EnclaveStatusInitializing means the enclave is being set up
	EnclaveStatusInitializing EnclaveStatus = "initializing"
	
	// EnclaveStatusRunning means the enclave is operational
	EnclaveStatusRunning EnclaveStatus = "running"
	
	// EnclaveStatusError means the enclave encountered an error
	EnclaveStatusError EnclaveStatus = "error"
	
	// EnclaveStatusClosed means the enclave has been shut down
	EnclaveStatusClosed EnclaveStatus = "closed"
)

// Enclave represents a trusted execution environment instance
type Enclave struct {
	ID           string        // Unique identifier for the enclave
	Provider     string        // TEE provider (azure, aws, etc.)
	Status       EnclaveStatus // Current status of the enclave
	CreatedAt    time.Time     // When the enclave was created
	LastActive   time.Time     // When the enclave was last used
	Error        string        // Last error message if status is error
	MemoryUsage  int64         // Current memory usage in bytes
	MemoryLimit  int64         // Maximum memory allowed in bytes
	Attestation  []byte        // Attestation report, if available
	AttestationExpiry time.Time // When the attestation expires
	
	// Internal state management
	mu          sync.RWMutex
	initialized bool
}

// NewEnclave creates a new enclave instance
func NewEnclave(id string, provider string, memoryLimit int64) *Enclave {
	now := time.Now()
	return &Enclave{
		ID:          id,
		Provider:    provider,
		Status:      EnclaveStatusInitializing,
		CreatedAt:   now,
		LastActive:  now,
		MemoryLimit: memoryLimit,
	}
}

// SetStatus updates the enclave status
func (e *Enclave) SetStatus(status EnclaveStatus, errorMsg string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.Status = status
	if status == EnclaveStatusError {
		e.Error = errorMsg
	}
	
	if status == EnclaveStatusRunning {
		e.initialized = true
	}
}

// SetAttestation sets the attestation report for the enclave
func (e *Enclave) SetAttestation(attestation []byte, expiry time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.Attestation = attestation
	e.AttestationExpiry = expiry
}

// UpdateMemoryUsage updates the current memory usage of the enclave
func (e *Enclave) UpdateMemoryUsage(usage int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.MemoryUsage = usage
	e.LastActive = time.Now()
}

// IsAttestationValid checks if the attestation is still valid
func (e *Enclave) IsAttestationValid() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	// Attestation is valid if it exists and hasn't expired
	return len(e.Attestation) > 0 && time.Now().Before(e.AttestationExpiry)
}

// IsInitialized checks if the enclave is initialized and ready for use
func (e *Enclave) IsInitialized() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	return e.initialized
}

// Close marks the enclave as closed
func (e *Enclave) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.Status = EnclaveStatusClosed
	e.LastActive = time.Now()
}
