package middleware

import (
	"net/http"

	"github.com/nprimo/quick/web"
)

// CORS is a middleware that adds Cross-Origin Resource Sharing headers.
func CORS(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) Middleware {
	return func(next web.HandlerFuncWithError) web.HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return next(w, r)
			}

			// Check if origin is allowed
			isOriginAllowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					isOriginAllowed = true
					break
				}
			}

			if !isOriginAllowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return nil
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", joinStrings(allowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", joinStrings(allowedHeaders, ","))
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
				w.WriteHeader(http.StatusOK)
				return nil
			}

			return next(w, r)
		}
	}
}

func joinStrings(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += sep + s[i]
	}
	return result
}
