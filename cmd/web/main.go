package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	homeHandler "github.com/AlejoGarat/snippetbox/internal/home/handlers"
	homeRepo "github.com/AlejoGarat/snippetbox/internal/home/repository"
	homeRoutes "github.com/AlejoGarat/snippetbox/internal/home/routes"
	snippetsHandler "github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
	snippetRepo "github.com/AlejoGarat/snippetbox/internal/snippets/repository"
	snippetsRoutes "github.com/AlejoGarat/snippetbox/internal/snippets/routes"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	snippetsRoutes.MakeRoutes(mux, snippetsHandler.New(
		errorLog,
		infoLog,
		&snippetRepo.SnippetRepo{DB: db},
	))
	homeRoutes.MakeRoutes(mux, homeHandler.New(
		errorLog,
		infoLog,
		&homeRepo.HomeRepo{DB: db},
	))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
