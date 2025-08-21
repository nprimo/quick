package middleware

import (
	"net/http"
	"time"

	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/web"
)

func Session(store sessions.Store) Middleware {
	return func(next web.HandlerFuncWithError) web.HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				return next(w, r)
			}

			session, err := store.Get(r.Context(), cookie.Value)
			if err != nil {
				return next(w, r)
			}

			if session.ExpiresAt.Before(time.Now()) {
				_ = store.Delete(r.Context(), session.ID)
				return next(w, r)
			}

			ctx := sessions.WithUserID(r.Context(), session.UserID)
			return next(w, r.WithContext(ctx))
		}
	}
}
