package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nprimo/quick/web"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logger(log *slog.Logger) Middleware {
	return func(next web.HandlerFuncWithError) web.HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now().UTC()

			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			err := next(recorder, r)

			log.Info("request",
				"method", r.Method,
				"status", recorder.status,
				"URI", r.RequestURI,
				"duration", time.Since(start),
				"remote_addr", r.RemoteAddr,
			)
			return err
		}
	}
}
