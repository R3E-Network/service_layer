package functions

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/R3E-Network/service_layer/internal/tee"
)

// Executor manages JavaScript function execution
type Executor struct {
	teeManager *tee.Manager
	functions  map[string]*models.Function
	mu         sync.RWMutex
}

// NewExecutor creates a new function executor
func NewExecutor(teeManager *tee.Manager) *Executor {
	return &Executor{
		teeManager: teeManager,
		functions:  make(map[string]*models.Function),
	}
}

// RegisterFunction adds a function to the registry
func (e *Executor) RegisterFunction(function *models.Function) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.functions[function.ID]; exists {
		return fmt.Errorf("function with ID %s already exists", function.ID)
	}

	// Register with TEE manager
	if err := e.teeManager.RegisterFunction(function.ID, function.Code, function.SecretRefs); err != nil {
		return err
	}

	e.functions[function.ID] = function
	log.Printf("Function registered: %s", function.ID)

	return nil
}

// UpdateFunction modifies an existing function
func (e *Executor) UpdateFunction(function *models.Function) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.functions[function.ID]; !exists {
		return fmt.Errorf("function with ID %s does not exist", function.ID)
	}

	// Update in TEE manager
	if err := e.teeManager.RegisterFunction(function.ID, function.Code, function.SecretRefs); err != nil {
		return err
	}

	e.functions[function.ID] = function
	log.Printf("Function updated: %s", function.ID)

	return nil
}

// DeleteFunction removes a function
func (e *Executor) DeleteFunction(functionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.functions[functionID]; !exists {
		return fmt.Errorf("function with ID %s does not exist", functionID)
	}

	delete(e.functions, functionID)
	log.Printf("Function deleted: %s", functionID)

	return nil
}

// GetFunction retrieves a function

// ExecuteFunction runs a JavaScript function in the TEE
func (e *Executor) ExecuteFunction(ctx context.Context, functionID string, params map[string]interface{}) (*models.ExecutionResult, error) {
	e.mu.RLock()
	function, exists := e.functions[functionID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("function with ID %s does not exist", functionID)
	}

	startTime := time.Now()

	// Execute in TEE
	result, err := e.teeManager.ExecuteFunction(ctx, functionID, params)
	if err != nil {
		return nil, err
	}

	executionTime := time.Since(startTime)

	// Format result
	executionResult := &models.ExecutionResult{
		FunctionID:    functionID,
		Status:        "success",
		Result:        result,
		ExecutionTime: executionTime,
		Timestamp:     time.Now(),
	}

	log.Printf("Function %s executed in %v", functionID, executionTime)

	return executionResult, nil
}
