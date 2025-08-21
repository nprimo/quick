package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/middleware"
	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/users"
	"github.com/nprimo/quick/web"
)

func Router(
	itemHandler items.Handler,
	userHandler users.Handler,
	sessionStore sessions.Store,
	log *slog.Logger,
) http.Handler {

	mux := web.NewServerMuxError(log)

	// Sensible defaults for CORS. Customize origins for your specific needs.
	allowedOrigins := []string{
		fmt.Sprintf("http://localhost:%s", LISTENING_PORT),
		fmt.Sprintf("http://127.0.0.1:%s", LISTENING_PORT),
	}
	allowedMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
	allowedHeaders := []string{"Content-Type", "Authorization", "X-CSRF-Token"}

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
	mux.HandleFuncWithError("/", basicChain.Then(indexHandler))
	mux.HandleFuncWithError("GET /register", basicChain.Then(userHandler.Register))
	mux.HandleFuncWithError("POST /register", basicChain.Then(userHandler.RegisterPost))
	mux.HandleFuncWithError("GET /login", basicChain.Then(userHandler.Login))
	mux.HandleFuncWithError("POST /login", basicChain.Then(userHandler.LoginPost))
	mux.HandleFuncWithError("GET /panic", basicChain.Then(func(w http.ResponseWriter, r *http.Request) error {
		panic("testing recover")
	}))

	// Protected routes
	mux.HandleFuncWithError("GET /items", protectedChain.Then(itemHandler.ListItems))
	mux.HandleFuncWithError("GET /items/{id}", protectedChain.Then(itemHandler.GetItem))
	mux.HandleFuncWithError("GET /items/new", protectedChain.Then(itemHandler.AddItem))
	mux.HandleFuncWithError("POST /items/new", protectedChain.Then(itemHandler.AddItemPost))
	mux.HandleFuncWithError("GET /items/{id}/update", protectedChain.Then(itemHandler.UpdateItem))
	mux.HandleFuncWithError("POST /items/{id}", protectedChain.Then(itemHandler.UpdateItemPost))
	mux.HandleFuncWithError("POST /items/{id}/delete", protectedChain.Then(itemHandler.DeleteItem))
	mux.HandleFuncWithError("GET /logout", protectedChain.Then(userHandler.Logout))

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("ciao mamma")
	// return ui.Index().Render(r.Context(), w)
}
