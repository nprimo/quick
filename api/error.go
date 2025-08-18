package api

import (
	"context"
	"log/slog"
	"net/http"
)

type HTTPError struct {
	Internal error
	Message  string
	Code     int
}

func (e HTTPError) Error() string {
	return e.Message
}

func HandleError(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
	code int,
	log *slog.Logger,
) {
	// TODO: make this better - add error stack, etc...
	log.Error("error",
		"err", err)
	if err := ErrorView(http.StatusText(code)).Render(ctx, w); err != nil {
		panic(err)
	}
}
