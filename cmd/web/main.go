package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/users"
)

const LISTENING_PORT = "4321"

func main() {
	dbConn, err := sql.Open("sqlite3", "dev.db")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	log := slog.Default()

	itemsStore := items.NewDBStore(dbConn)
	itemsHandler := items.NewHandler(itemsStore, log)

	sessionsStore := sessions.NewDBStore(dbConn)
	usersStore := users.NewDBStore(dbConn)
	usersHandler := users.NewHandler(usersStore, sessionsStore, log)

	server := http.Server{
		// TODO: make this come from config
		Addr:         ":" + LISTENING_PORT,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      Router(itemsHandler, usersHandler, sessionsStore, log),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to listen and serve",
			"error", err)
	}
}
