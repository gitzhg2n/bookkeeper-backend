package routes

import (
	"net/http"
	"strconv"
	"strings"
)

func validateRequired(fields map[string]string) (string, bool) {
	for k, v := range fields {
		if strings.TrimSpace(v) == "" {
			return k, false
		}
	}
	return "", true
}

func parseUintParam(param string, paramName string, w http.ResponseWriter) (uint, bool) {
	if param == "" {
		writeJSONError(w, "missing "+paramName, http.StatusBadRequest)
		return 0, false
	}
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		writeJSONError(w, "invalid "+paramName, http.StatusBadRequest)
		return 0, false
	}
	return uint(id), true
}

func sanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}