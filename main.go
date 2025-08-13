package main

import (
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nprimo/quick/db"
	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/web"
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

	itemsStore := items.NewDBStore(dbConn)
	itemsHandler := items.NewHandler(itemsStore)
	log := slog.New(&slog.JSONHandler{})

	http.ListenAndServe(":4321", web.Router(itemsHandler, log))
}
