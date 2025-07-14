package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ciao mamma`))
	})

	http.ListenAndServe(":4321", nil)
}
