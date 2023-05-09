package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	fileServerRoutes "github.com/AlejoGarat/snippetbox/internal/fileserver/routes"
	homeHandler "github.com/AlejoGarat/snippetbox/internal/home/handlers"
	homeRepo "github.com/AlejoGarat/snippetbox/internal/home/repository"
	homeRoutes "github.com/AlejoGarat/snippetbox/internal/home/routes"
	homeService "github.com/AlejoGarat/snippetbox/internal/home/service"
	snippetsHandler "github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
	snippetRepo "github.com/AlejoGarat/snippetbox/internal/snippets/repository"
	snippetsRoutes "github.com/AlejoGarat/snippetbox/internal/snippets/routes"
	snippetService "github.com/AlejoGarat/snippetbox/internal/snippets/service"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
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

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httphelpers.NotFound(w)
	})

	fileServerRoutes.MakeRoutes(router)

	snippetsRoutes.MakeRoutes(router, snippetsHandler.New(
		errorLog,
		infoLog,
		snippetService.NewSnippetService(snippetRepo.NewSnippetRepo(db)),
	))

	homeRoutes.MakeRoutes(router, homeHandler.New(
		errorLog,
		infoLog,
		homeService.NewHomeService(homeRepo.NewSnippetRepo(db)),
	))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  router,
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
