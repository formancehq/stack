package http

import (
	"net/http"
)

func IsHTTPRequestSuccess(statusCode int) bool {
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return true
	}
	return false
}
