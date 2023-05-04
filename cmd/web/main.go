package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	commonHandler "github.com/AlejoGarat/snippetbox/internal/common/handlers"
	homeRoutes "github.com/AlejoGarat/snippetbox/internal/common/routes"
	snippetsHandler "github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
	snippetsRoutes "github.com/AlejoGarat/snippetbox/internal/snippets/routes"
)

func main() {
	addr := os.Getenv("SNIPPETBOX_ADDR")

	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	snippetsRoutes.MakeRoutes(mux, &snippetsHandler.Handler{})
	homeRoutes.MakeRoutes(mux, &commonHandler.Handler{})

	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
