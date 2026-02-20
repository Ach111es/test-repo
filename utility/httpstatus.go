package utility

import (
	"net/http"
	"strings"
)

func ErrorHttpStatus(err string) (int, string) {
	defaultStatusMessage := "Bad Request"

	err = strings.ToLower(err)

	errorStatusMap := map[string]struct {
		Code   int
		Status string
	}{
		"duplicate":          {http.StatusBadRequest, defaultStatusMessage},
		"parsing time":       {http.StatusBadRequest, defaultStatusMessage},
		"required":           {http.StatusBadRequest, defaultStatusMessage},
		"invalid":            {http.StatusBadRequest, defaultStatusMessage},
		"not found":          {http.StatusNotFound, "Not Found"},
		"unauthorized":       {http.StatusUnauthorized, "Unauthorized"},
		"forbidden":          {http.StatusForbidden, "Forbidden"},
		"connection refused": {http.StatusInternalServerError, "Internal Server Error"},
	}

	for key, status := range errorStatusMap {
		if strings.Contains(err, key) {
			return status.Code, status.Status
		}
	}

	return http.StatusBadRequest, defaultStatusMessage
}
