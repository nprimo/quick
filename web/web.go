package web

import (
	"log/slog"
	"net/http"

	"github.com/nprimo/quick/ui"
)

type Error struct {
	Code          int
	Message       string
	InternalError error
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, err error, message string) Error {
	if len(message) == 0 {
		message = http.StatusText(code)
	}
	return Error{
		Code:          code,
		InternalError: err,
		Message:       message,
	}
}

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
			s.log.Error("error in handler",
				"err", err,
				"method", r.Method,
				"uri", r.RequestURI,
				"pattern", pattern,
			)
			w.WriteHeader(http.StatusInternalServerError)
			ui.Error(err).Render(r.Context(), w)
		}
	})
}
