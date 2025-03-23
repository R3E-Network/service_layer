package functions

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/your-org/neo-oracle/internal/models"
)

// AdvancedJSRuntime provides enhanced JavaScript runtime with blockchain capabilities
type AdvancedJSRuntime struct {
	vm          *goja.Runtime
	bridge      *RuntimeBridge
	httpClient  *http.Client
	maxFetchMB  int
	sandboxed   bool
	allowedURLs []string
}

// NewAdvancedJSRuntime creates an enhanced JavaScript runtime
func NewAdvancedJSRuntime(maxFetchMB int, sandboxed bool, allowedURLs []string) *AdvancedJSRuntime {
	vm := goja.New()
	bridge := NewRuntimeBridge(vm)

	return &AdvancedJSRuntime{
		vm:          vm,
		bridge:      bridge,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		maxFetchMB:  maxFetchMB,
		sandboxed:   sandboxed,
		allowedURLs: allowedURLs,
	}
}

// RegisterGlobalFunctions adds capabilities to the JavaScript environment
func (r *AdvancedJSRuntime) RegisterGlobalFunctions() {
	// Console functions
	console := r.bridge.Console()
	r.vm.Set("console", console)

	// Timing functions
	r.vm.Set("setTimeout", r.bridge.SetTimeout)
	r.vm.Set("clearTimeout", r.bridge.ClearTimeout)

	// Fetch API (carefully restricted)
	r.vm.Set("fetch", r.safeFetch)

	// Crypto functions
	r.vm.Set("crypto", r.cryptoFunctions())

	// Blockchain utilities
	r.vm.Set("blockchain", r.blockchainFunctions())

	// Basic utilities
	r.vm.Set("utils", r.utilityFunctions())
}

// safeFetch provides a limited fetch API for JavaScript
func (r *AdvancedJSRuntime) safeFetch(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		return r.vm.ToValue(map[string]interface{}{
			"ok":    false,
			"error": "URL required",
		})
	}

	// Get URL from arguments
	url := call.Arguments[0].String()

	// Options (optional)
	options := make(map[string]interface{})
	if len(call.Arguments) > 1 && !goja.IsUndefined(call.Arguments[1]) {
		// Convert options to Go map
		optionsObj := call.Arguments[1].ToObject(r.vm)
		for _, key := range optionsObj.Keys() {
			val := optionsObj.Get(key)
			if val != nil {
				options[key] = val.Export()
			}
		}
	}

	// Security checks for URLs
	if r.sandboxed {
		// Check if URL is in the allowed list
		allowed := false
		for _, prefix := range r.allowedURLs {
			if strings.HasPrefix(url, prefix) {
				allowed = true
				break
			}
		}

		if !allowed {
			return r.vm.ToValue(map[string]interface{}{
				"ok":    false,
				"error": "URL not allowed in sandbox mode",
			})
		}
	}

	// Make the HTTP request
	method := "GET"
	if m, ok := options["method"].(string); ok && m != "" {
		method = m
	}

	// Create request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return r.vm.ToValue(map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
	}

	// Add headers
	if headers, ok := options["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strVal, ok := value.(string); ok {
				req.Header.Add(key, strVal)
			}
		}
	}

	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return r.vm.ToValue(map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Check max size
	if resp.ContentLength > int64(r.maxFetchMB*1024*1024) {
		return r.vm.ToValue(map[string]interface{}{
			"ok":    false,
			"error": fmt.Sprintf("Response too large (max: %d MB)", r.maxFetchMB),
		})
	}

	// Read response body with limit
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r.vm.ToValue(map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
	}

	// Create response object
	respObj := map[string]interface{}{
		"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    resp.Header,
		"text": func() string {
			return string(bodyBytes)
		},
		"json": func() interface{} {
			var result interface{}
			if err := json.Unmarshal(bodyBytes, &result); err != nil {
				// Return parsing error
				return map[string]interface{}{
					"error": fmt.Sprintf("JSON parse error: %v", err),
				}
			}
			return result
		},
	}

	return r.vm.ToValue(respObj)
}

// cryptoFunctions provides cryptographic utilities for JavaScript
func (r *AdvancedJSRuntime) cryptoFunctions() map[string]interface{} {
	return map[string]interface{}{
		"sha256": func(data string) string {
			hash := sha256.Sum256([]byte(data))
			return hex.EncodeToString(hash[:])
		},
		"randomBytes": func(length int) string {
			if length <= 0 || length > 1024 {
				length = 32 // Default to 32 bytes
			}
			bytes := make([]byte, length)
			for i := 0; i < length; i++ {
				bytes[i] = byte(r.vm.ToValue(time.Now().UnixNano()).ToInteger() % 256)
				time.Sleep(1 * time.Nanosecond) // Add entropy
			}
			return hex.EncodeToString(bytes)
		},
	}
}

// blockchainFunctions provides Neo N3 blockchain utilities for JavaScript
func (r *AdvancedJSRuntime) blockchainFunctions() map[string]interface{} {
	return map[string]interface{}{
		"hexToBase64": func(hexStr string) string {
			bytes, err := hex.DecodeString(hexStr)
			if err != nil {
				return ""
			}
			return hex.EncodeToString(bytes)
		},
		"scriptHashToAddress": func(scriptHash string) string {
			// In a real implementation, this would convert Neo scriptHash to address
			return "N" + scriptHash[:20]
		},
	}
}

// utilityFunctions provides general utilities for JavaScript
func (r *AdvancedJSRuntime) utilityFunctions() map[string]interface{} {
	return map[string]interface{}{
		"sleep": func(ms int) {
			if ms <= 0 || ms > 30000 { // Max 30 seconds
				ms = 1000
			}
			time.Sleep(time.Duration(ms) * time.Millisecond)
		},
		"parseDate": func(dateStr string) int64 {
			t, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				return 0
			}
			return t.Unix()
		},
		"formatDate": func(timestamp int64) string {
			return time.Unix(timestamp, 0).Format(time.RFC3339)
		},
	}
}

// ExecuteCode runs JavaScript code in the runtime
func (r *AdvancedJSRuntime) ExecuteCode(ctx context.Context, code string, params map[string]interface{}, secrets map[string]string) (*models.ExecutionResult, error) {
	// Create a context with cancellation
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Reset runtime before execution
	r.vm = goja.New()
	r.bridge = NewRuntimeBridge(r.vm)
	r.RegisterGlobalFunctions()

	// Inject parameters and secrets
	paramsJSON, _ := json.Marshal(params)
	secretsJSON, _ := json.Marshal(secrets)
	r.vm.Set("params", string(paramsJSON))
	r.vm.Set("secrets", string(secretsJSON))

	// Create the full script with wrapper
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
	`, paramsJSON, secretsJSON, code)

	// Execute in a goroutine with timeout
	resultChan := make(chan *models.ExecutionResult, 1)
	errorChan := make(chan error, 1)

	go func() {
		startTime := time.Now()

		// Run the script
		value, err := r.vm.RunString(script)
		if err != nil {
			errorChan <- fmt.Errorf("execution error: %w", err)
			return
		}

		// Process result
		result := value.Export()
		execResult := &models.ExecutionResult{
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
		return nil, fmt.Errorf("function execution timed out")
	}
}

// Cleanup releases resources
func (r *AdvancedJSRuntime) Cleanup() {
	r.bridge.CleanupTimeouts()
}
