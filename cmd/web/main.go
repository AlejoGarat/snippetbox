package main

import (
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

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	snippetsRoutes.MakeRoutes(mux, &snippetsHandler.Handler{})
	homeRoutes.MakeRoutes(mux, &commonHandler.Handler{})

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
