package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nprimo/quick/db"
	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/users"
)

func main() {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	if err := db.Init(dbConn); err != nil {
		panic(err)
	}
	// TODO: make this only if flag enabled (for example)
	if err := db.Seed(dbConn); err != nil {
		panic(err)
	}

	log := slog.Default()

	itemsStore := items.NewDBStore(dbConn)
	itemsHandler := items.NewHandler(itemsStore, log)

	usersStore := users.NewDBStore(dbConn)
	usersHandler := users.NewHandler(usersStore, log)

	server := http.Server{
		// TODO: make this come from config
		Addr:         ":4321",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      Router(itemsHandler, usersHandler, log),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to listen and serve",
			"error", err)
	}
}
