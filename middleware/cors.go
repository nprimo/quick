package middleware

import "net/http"

// CORS is a middleware that adds Cross-Origin Resource Sharing headers.
func CORS(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
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
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", joinStrings(allowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", joinStrings(allowedHeaders, ","))
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
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
