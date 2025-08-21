package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"slices"

	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/web"
)

const (
	csrfHeader  = "X-CSRF-Token"
	csrfFormKey = "csrf_token"
)

// CSRF is a middleware that provides Cross-Site Request Forgery protection.
func CSRF(store sessions.Store) Middleware {
	return func(next web.HandlerFuncWithError) web.HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			// Get session ID from cookie
			cookie, err := r.Cookie("session_id")
			var sessionID string
			if err == nil {
				sessionID = cookie.Value
			}

			// Get session from store
			session, err := store.Get(r.Context(), sessionID)
			if err != nil {
				// If session not found or error, create a dummy session for token generation
				session = sessions.Session{}
			}

			// Check if token exists in session
			token := session.CSRFToken
			if token == "" {
				// Generate new token if not present
				token, err = generateCSRFToken()
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return nil
				}
				// Update session with new token
				session.CSRFToken = token
				// Save the session (even if it's a new one)
				if sessionID != "" {
					if err = store.Save(r.Context(), session); err != nil {
						http.Error(w, "failed to save session", http.StatusForbidden)
						return nil
					}
				}
			}

			// Validate token for POST requests
			methods := []string{
				http.MethodDelete,
				http.MethodPost,
			}
			if slices.Contains(methods, r.Method) {
				submittedToken := r.Header.Get(csrfHeader)
				if submittedToken == "" {
					submittedToken = r.FormValue(csrfFormKey)
				}

				if submittedToken != token {
					http.Error(w, "CSRF token mismatch", http.StatusForbidden)
					return nil
				}
			}

			// Add token to context for templates
			ctx := sessions.WithCSRFToken(r.Context(), token)
			return next(w, r.WithContext(ctx))
		}
	}
}

func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
