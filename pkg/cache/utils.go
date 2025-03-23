package cache

import (
	"fmt"
)

// FormatFunctionDetailsKey formats a key for function details
func FormatFunctionDetailsKey(id string) string {
	return fmt.Sprintf("function:details:%s", id)
}

// FormatFunctionListKey formats a key for a list of functions with pagination
func FormatFunctionListKey(userID int, page, limit int) string {
	return fmt.Sprintf("function:list:%d:%d:%d", userID, page, limit)
}

// BuildInvalidationPatterns builds patterns for cache invalidation
func BuildInvalidationPatterns(objectType string, id string, userID int) []string {
	patterns := []string{}
	
	switch objectType {
	case "function":
		// Invalidate specific function
		patterns = append(patterns, fmt.Sprintf("function:details:%s", id))
		
		// Invalidate function lists for this user
		if userID > 0 {
			patterns = append(patterns, fmt.Sprintf("function:list:%d:*", userID))
		}
	case "user":
		// Invalidate user details
		patterns = append(patterns, fmt.Sprintf("user:details:%s", id))
		
		// Invalidate user lists
		patterns = append(patterns, "user:list:*")
	}
	
	return patterns
}
