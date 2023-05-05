package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	homeHandler "github.com/AlejoGarat/snippetbox/internal/home/handlers"
	homeRoutes "github.com/AlejoGarat/snippetbox/internal/home/routes"
	snippetsHandler "github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
	snippetsRoutes "github.com/AlejoGarat/snippetbox/internal/snippets/routes"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

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

	snippetsRoutes.MakeRoutes(mux, &snippetsHandler.Handler{
		ErrorLog: errorLog, InfoLog: infoLog,
	})
	homeRoutes.MakeRoutes(mux, &homeHandler.Handler{
		ErrorLog: errorLog, InfoLog: infoLog,
	})

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
