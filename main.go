package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nprimo/quick/items"
	"github.com/nprimo/quick/web"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// TODO: make this a migration
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, name TEXT, quantity INTEGER)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`insert into items(name, quantity) values('banana', 1)`)
	if err != nil {
		log.Fatal(err)
	}

	itemsStore := items.NewDBStore(db)
	itemsHandler := items.NewHandler(itemsStore)

	http.ListenAndServe(":4321", web.Router(itemsHandler))
}

