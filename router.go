package main

import (
	"log/slog"
	"net/http"

	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/middleware"
	"github.com/nprimo/quick/ui"
	"github.com/nprimo/quick/users"
	"github.com/nprimo/quick/web"
)

func Router(
	itemHandler items.Handler,
	userHandler users.Handler,
	log *slog.Logger,
) http.Handler {
	mux := web.NewErrorMux(log)

	mux.HandleErrorFunc("/", func(w http.ResponseWriter, r *http.Request) error {
		if err := ui.Index().Render(r.Context(), w); err != nil {
			return web.NewError(http.StatusInternalServerError, err, "")
		}
		return nil
	})

	mux.HandleErrorFunc("GET /items", itemHandler.ListItems)
	mux.HandleErrorFunc("GET /items/{id}", itemHandler.GetItem)

	mux.HandleErrorFunc("GET /items/{id}/update", itemHandler.UpdateItem)
	mux.HandleErrorFunc("POST /items/{id}", itemHandler.UpdateItemPost)

	mux.HandleErrorFunc("GET /items/{id}/delete", itemHandler.DeleteItem)

	mux.HandleErrorFunc("GET /items/new", itemHandler.AddItem)
	mux.HandleErrorFunc("POST /items/new", itemHandler.AddItemPost)

	mux.HandleErrorFunc("GET /register", userHandler.Register)
	mux.HandleErrorFunc("POST /register", userHandler.RegisterPost)
	mux.HandleErrorFunc("GET /login", userHandler.Login)
	mux.HandleErrorFunc("POST /login", userHandler.LoginPost)

	wrapped := middleware.Logger(mux, log)
	return wrapped
}
