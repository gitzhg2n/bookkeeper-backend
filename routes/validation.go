package routes

import (
	"net/http"
	"strconv"
	"strings"
)

// validateRequired checks if all required fields are present and non-empty
func validateRequired(fields map[string]string) (string, bool) {
	for fieldName, value := range fields {
		if strings.TrimSpace(value) == "" {
			return fieldName, false
		}
	}
	return "", true
}

// parseUintParam safely parses a URL parameter to uint
func parseUintParam(param string, paramName string, w http.ResponseWriter) (uint, bool) {
	if param == "" {
		writeJSONError(w, "Missing "+paramName, http.StatusBadRequest)
		return 0, false
	}
	
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		writeJSONError(w, "Invalid "+paramName, http.StatusBadRequest)
		return 0, false
	}
	
	return uint(id), true
}

// sanitizeString performs basic sanitization on string input
func sanitizeString(input string) string {
	// Remove leading/trailing whitespace
	input = strings.TrimSpace(input)
	
	// Replace any null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	return input
}