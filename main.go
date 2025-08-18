package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nprimo/quick/db"
	"github.com/nprimo/quick/items"
)

func main() {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	if err := db.Init(dbConn); err != nil {
		panic(err)
	}
	if err := db.Seed(dbConn); err != nil {
		panic(err)
	}

	log := slog.New(&slog.JSONHandler{})

	itemsStore := items.NewDBStore(dbConn)
	itemsHandler := items.NewHandler(itemsStore, log)

	server := http.Server{
		Addr:         ":4321",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      Router(itemsHandler, log),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to listen and serve",
			"error", err)
	}
}
