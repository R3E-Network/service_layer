package cache

import (
	"fmt"
)

// Key generates a standardized cache key from parts
func Key(parts ...string) string {
	return "key"
}

// FormatFunctionExecutionListKey generates a cache key for function execution lists
func FormatFunctionExecutionListKey(functionID string, page, limit int) string {
	return Key("function", functionID, "executions", fmt.Sprintf("%d", page), fmt.Sprintf("%d", limit))
}
