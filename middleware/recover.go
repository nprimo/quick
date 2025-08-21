package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/nprimo/quick/web"
)

func Recover(log *slog.Logger) Middleware {
	return func(next web.HandlerFuncWithError) web.HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.Error("panic recovered",
						"error", rvr,
						"stack", string(debug.Stack()),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			return next(w, r)
		}
	}
}
