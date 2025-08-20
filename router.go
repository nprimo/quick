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

	// Sensible defaults for CORS. Customize origins for your specific needs.
	allowedOrigins := []string{
		fmt.Sprintf("http://localhost:%s", LISTENING_PORT),
		fmt.Sprintf("http://127.0.0.1:%s", LISTENING_PORT),
	}
	allowedMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
	allowedHeaders := []string{"Content-Type", "Authorization", "X-CSRF-Token"}

	errorHandler := web.ErrorHandler(log)
	basicChain := middleware.New(
		middleware.CORS(allowedOrigins, allowedMethods, allowedHeaders), // CORS first
		middleware.Logger(log),
		middleware.Session(sessionStore),
	)
	protectedChain := basicChain.Append(
		middleware.RequireLogin,
		middleware.CSRF(sessionStore),
	)

	// Public routes
	mux.Handle("/", basicChain.Then(errorHandler(indexHandler)))
	mux.Handle("GET /register", basicChain.Then(errorHandler(userHandler.Register)))
	mux.Handle("POST /register", basicChain.Then(errorHandler(userHandler.RegisterPost)))
	mux.Handle("GET /login", basicChain.Then(errorHandler(userHandler.Login)))
	mux.Handle("POST /login", basicChain.Then(errorHandler(userHandler.LoginPost)))

	// Protected routes
	mux.Handle("GET /items", protectedChain.Then(errorHandler(itemHandler.ListItems)))
	mux.Handle("GET /items/{id}", protectedChain.Then(errorHandler(itemHandler.GetItem)))
	mux.Handle("GET /items/new", protectedChain.Then(errorHandler(itemHandler.AddItem)))
	mux.Handle("POST /items/new", protectedChain.Then(errorHandler(itemHandler.AddItemPost)))
	mux.Handle("GET /items/{id}/update", protectedChain.Then(errorHandler(itemHandler.UpdateItem)))
	mux.Handle("POST /items/{id}", protectedChain.Then(errorHandler(itemHandler.UpdateItemPost)))
	mux.Handle("POST /items/{id}/delete", protectedChain.Then(errorHandler(itemHandler.DeleteItem)))
	mux.Handle("GET /logout", protectedChain.Then(errorHandler(userHandler.Logout)))

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	if err := ui.Index().Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	return nil
}
