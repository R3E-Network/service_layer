package functions

import (
	"fmt"
	"log"
	"time"

	"github.com/dop251/goja"
)

// RuntimeBridge provides utility functions to JavaScript code
type RuntimeBridge struct {
	vm       *goja.Runtime
	timeouts map[string]*time.Timer
}

// NewRuntimeBridge creates a new bridge for JavaScript
func NewRuntimeBridge(vm *goja.Runtime) *RuntimeBridge {
	return &RuntimeBridge{
		vm:       vm,
		timeouts: make(map[string]*time.Timer),
	}
}

// Console provides console logging for JavaScript
func (rb *RuntimeBridge) Console() map[string]interface{} {
	return map[string]interface{}{
		"log": func(args ...interface{}) {
			log.Printf("JS: %v", args)
		},
		"error": func(args ...interface{}) {
			log.Printf("JS ERROR: %v", args)
		},
		"warn": func(args ...interface{}) {
			log.Printf("JS WARNING: %v", args)
		},
	}
}

// SetTimeout implements JavaScript setTimeout
func (rb *RuntimeBridge) SetTimeout(callback goja.Value, ms goja.Value) goja.Value {
	if !callback.IsFunction() {
		return goja.Undefined()
	}

	// Get milliseconds
	msInt := int64(ms.ToInteger())
	if msInt < 0 {
		msInt = 0
	}

	// Generate timeout ID
	timeoutID := fmt.Sprintf("timeout-%d", time.Now().UnixNano())

	// Create timer
	timer := time.AfterFunc(time.Duration(msInt)*time.Millisecond, func() {
		_, _ = callback.Call(goja.Undefined())
		delete(rb.timeouts, timeoutID)
	})

	rb.timeouts[timeoutID] = timer
	return rb.vm.ToValue(timeoutID)
}

// ClearTimeout implements JavaScript clearTimeout
func (rb *RuntimeBridge) ClearTimeout(timeoutID goja.Value) {
	id := timeoutID.String()
	if timer, exists := rb.timeouts[id]; exists {
		timer.Stop()
		delete(rb.timeouts, id)
	}
}

// FetchImplementation provides a simplified fetch API for JavaScript
func (rb *RuntimeBridge) FetchImplementation(url string, options map[string]interface{}) map[string]interface{} {
	// In a real implementation, this would make HTTP requests
	// For now, we'll simulate responses

	log.Printf("JS fetch: %s", url)

	return map[string]interface{}{
		"ok":         true,
		"status":     200,
		"statusText": "OK",
		"json": func() map[string]interface{} {
			return map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"value":     "Simulated response",
					"timestamp": time.Now().Unix(),
				},
			}
		},
		"text": func() string {
			return "Simulated response text"
		},
	}
}

// CleanupTimeouts stops all pending timeouts
func (rb *RuntimeBridge) CleanupTimeouts() {
	for id, timer := range rb.timeouts {
		timer.Stop()
		delete(rb.timeouts, id)
	}
}
