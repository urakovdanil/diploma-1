package middleware

import (
	"diploma-1/internal/logger"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		crw, ok := w.(*customResponseWriter)
		if !ok {
			logger.Warn(r.Context(), "unable to cast response writer to customResponseWriter, request logging will not work")
		}
		if crw.statusCode == 0 {
			crw.statusCode = http.StatusOK
		}
		took := time.Since(start)
		if crw.statusCode >= http.StatusInternalServerError {
			logger.Errorf(r.Context(), "request: %v, %v, took: %v; reponse: %v, size: %v, details: %v", r.Method, r.URL.Path, took, crw.statusCode, crw.responseSize, crw.buf.String())
		} else {
			logger.Debugf(r.Context(), "request: %v, %v, took: %v; reponse: %v, size: %v", r.Method, r.URL.Path, took, crw.statusCode, crw.responseSize)
		}
	})
}
