package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/middleware"
	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/ui"
	"github.com/nprimo/quick/users"
	"github.com/nprimo/quick/web"
)

func Router(
	itemHandler items.Handler,
	userHandler users.Handler,
	sessionStore sessions.Store,
	log *slog.Logger,
) http.Handler {

	mux := http.NewServeMux()

	must := func(h http.Handler, err error) http.Handler {
		if err != nil {
			panic(err)
		}
		return h
	}

	// Sensible defaults for CORS. Customize origins for your specific needs.
	allowedOrigins := []string{
		fmt.Sprintf("http://localhost:%s", LISTENING_PORT),
		fmt.Sprintf("http://127.0.0.1:%s", LISTENING_PORT),
	}
	allowedMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
	allowedHeaders := []string{"Content-Type", "Authorization", "X-CSRF-Token"}

	errorHandler := web.ErrorHandler(log)
	basicChain := middleware.New(
		middleware.Recover(log), // Recover from panics first
		middleware.CORS(allowedOrigins, allowedMethods, allowedHeaders), // Then CORS
		middleware.Logger(log),
		middleware.Session(sessionStore),
	)
	protectedChain := basicChain.Append(
		middleware.RequireLogin,
		middleware.CSRF(sessionStore),
	)

	// Public routes
	mux.Handle("/", must(basicChain.Then(errorHandler(indexHandler))))
	mux.Handle("GET /register", must(basicChain.Then(errorHandler(userHandler.Register))))
	mux.Handle("POST /register", must(basicChain.Then(errorHandler(userHandler.RegisterPost))))
	mux.Handle("GET /login", must(basicChain.Then(errorHandler(userHandler.Login))))
	mux.Handle("POST /login", must(basicChain.Then(errorHandler(userHandler.LoginPost))))
	mux.Handle("GET /panic", must(basicChain.Then(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { panic("testing recover") }),
	)))

	// Protected routes
	mux.Handle("GET /items", must(protectedChain.Then(errorHandler(itemHandler.ListItems))))
	mux.Handle("GET /items/{id}", must(protectedChain.Then(errorHandler(itemHandler.GetItem))))
	mux.Handle("GET /items/new", must(protectedChain.Then(errorHandler(itemHandler.AddItem))))
	mux.Handle("POST /items/new", must(protectedChain.Then(errorHandler(itemHandler.AddItemPost))))
	mux.Handle("GET /items/{id}/update", must(protectedChain.Then(errorHandler(itemHandler.UpdateItem))))
	mux.Handle("POST /items/{id}", must(protectedChain.Then(errorHandler(itemHandler.UpdateItemPost))))
	mux.Handle("POST /items/{id}/delete", must(protectedChain.Then(errorHandler(itemHandler.DeleteItem))))
	mux.Handle("GET /logout", must(protectedChain.Then(errorHandler(userHandler.Logout))))

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	if err := ui.Index().Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	return nil
}
