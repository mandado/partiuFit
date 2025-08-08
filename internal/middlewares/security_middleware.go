package middlewares

import (
	"go.uber.org/zap"
	"net/http"
)

type SecurityMiddleware struct {
	Logger *zap.SugaredLogger
}

func NewSecurityMiddleware(logger *zap.SugaredLogger) *SecurityMiddleware {
	return &SecurityMiddleware{
		Logger: logger,
	}
}

// SecurityHeaders adds security headers appropriate for JSON API
func (sm *SecurityMiddleware) SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing (IMPORTANTE para APIs)
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Strict Transport Security (HTTPS only) - IMPORTANTE para APIs
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Remove server header for security (esconde versão do servidor)
		w.Header().Set("Server", "")

		// Referrer Policy - útil para APIs também
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Para APIs, definir Content-Type explícito
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		// Cache control para endpoints sensíveis (não cachear responses com dados pessoais)
		if r.Header.Get("Authorization") != "" {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		next.ServeHTTP(w, r)
	})
}
