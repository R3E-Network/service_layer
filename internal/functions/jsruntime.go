package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dop251/goja"
	"github.com/your-org/neo-oracle/internal/models"
	"github.com/your-org/neo-oracle/internal/secrets"
)

// JSRuntime provides JavaScript execution capabilities
type JSRuntime struct {
	secretsSvc *secrets.Service
	timeout    time.Duration
}

// NewJSRuntime creates a new JavaScript runtime
func NewJSRuntime(secretsSvc *secrets.Service, timeoutSec int) *JSRuntime {
	return &JSRuntime{
		secretsSvc: secretsSvc,
		timeout:    time.Duration(timeoutSec) * time.Second,
	}
}

// ExecuteFunction runs a JavaScript function with the given parameters
func (r *JSRuntime) ExecuteFunction(ctx context.Context, function *models.Function, params map[string]interface{}) (*models.ExecutionResult, error) {
	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Create JavaScript VM
	vm := goja.New()

	// Add console.log implementation
	console := map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			args := make([]interface{}, len(call.Arguments))
			for i, arg := range call.Arguments {
				args[i] = arg.String()
			}
			log.Printf("JS Console: %v", args)
			return goja.Undefined()
		},
	}
	vm.Set("console", console)

	// Add setTimeout implementation (simplified)
	vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return goja.Undefined()
		}
		return goja.Undefined()
	})

	// Add fetch implementation (simplified)
	vm.Set("fetch", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(map[string]interface{}{
				"error": "URL required",
			})
		}

		// In real implementation, this would make actual HTTP requests
		// For security, this would be carefully restricted

		return vm.ToValue(map[string]interface{}{
			"json": func() map[string]interface{} {
				return map[string]interface{}{
					"result": "Simulated fetch response",
				}
			},
			"text": func() string {
				return "Simulated fetch response text"
			},
		})
	})

	// Inject parameters
	paramsJSON, _ := json.Marshal(params)
	vm.Set("params", string(paramsJSON))

	// Inject secrets if available
	secrets := make(map[string]string)
	for _, secretRef := range function.SecretRefs {
		secretValue, err := r.secretsSvc.RetrieveSecret(secretRef)
		if err != nil {
			log.Printf("Warning: could not retrieve secret %s: %v", secretRef, err)
			continue
		}
		secrets[secretRef] = secretValue
	}
	secretsJSON, _ := json.Marshal(secrets)
	vm.Set("secrets", string(secretsJSON))

	// Prepare the full script with wrapper
	script := fmt.Sprintf(`
		const params = JSON.parse(%s);
		const secrets = JSON.parse(%s);
		
		// User function
		%s
		
		// Execute and return result
		(function() {
			try {
				const result = main(params, secrets);
				return { success: true, result: result };
			} catch (error) {
				return { 
					success: false, 
					error: error.toString(),
					stack: error.stack
				};
			}
		})();
	`, paramsJSON, secretsJSON, function.Code)

	// Run the script with timeout
	resultChan := make(chan *models.ExecutionResult, 1)
	errorChan := make(chan error, 1)

	go func() {
		startTime := time.Now()

		// Execute JavaScript
		value, err := vm.RunString(script)
		if err != nil {
			errorChan <- fmt.Errorf("execution error: %w", err)
			return
		}

		// Extract result
		result := value.Export()
		execResult := &models.ExecutionResult{
			FunctionID:    function.ID,
			ExecutionTime: time.Since(startTime),
			Timestamp:     time.Now(),
		}

		// Check for success
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			errorChan <- fmt.Errorf("invalid result format")
			return
		}

		success, _ := resultMap["success"].(bool)
		if !success {
			execResult.Status = "error"
			execResult.Error = fmt.Sprintf("%v", resultMap["error"])
		} else {
			execResult.Status = "success"
			execResult.Result = resultMap["result"]
		}

		resultChan <- execResult
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-execCtx.Done():
		return nil, fmt.Errorf("function execution timed out after %v", r.timeout)
	}
}
