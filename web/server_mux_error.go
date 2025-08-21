package web

import (
	"log/slog"
	"net/http"

	"github.com/nprimo/quick/ui"
)

type HandlerFuncWithError func(w http.ResponseWriter, r *http.Request) error

type ServerMuxError struct {
	*http.ServeMux
	log *slog.Logger
}

func NewServerMuxError(log *slog.Logger) *ServerMuxError {
	return &ServerMuxError{
		ServeMux: http.NewServeMux(),
		log:      log,
	}
}

func (s *ServerMuxError) HandleFuncWithError(pattern string, handler HandlerFuncWithError) {
	s.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			s.log.Error("error in handler", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			ui.Error(err).Render(r.Context(), w)
		}
	})
}
