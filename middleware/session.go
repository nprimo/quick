package middleware

import (
	"net/http"
	"time"

	"github.com/nprimo/quick/sessions"
)

func Session(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			session, err := store.Get(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if session.ExpiresAt.Before(time.Now()) {
				_ = store.Delete(r.Context(), session.ID)
				next.ServeHTTP(w, r)
				return
			}

			ctx := sessions.WithUserID(r.Context(), session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
