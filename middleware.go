package io

import "net/http"

// MiddlewareLog writes a log message for
func MiddlewareLog(next http.HandlerFunc, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		logger.log(severityLow, "Log message", r)
	}
}
