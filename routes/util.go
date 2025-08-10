package routes

import (
	"strconv"
)

// parseUintString parses a string to uint, returning the value and whether parsing succeeded
func parseUintString(s string) (uint, bool) {
	id64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(id64), true
}