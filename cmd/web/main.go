package main

import (
	"log"
	"net/http"

	snippetsHandler "github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", snippetsHandler.Home)
	mux.HandleFunc("/snippet/view", snippetsHandler.SnippetView)
	mux.HandleFunc("/snippet/create", snippetsHandler.SnippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
