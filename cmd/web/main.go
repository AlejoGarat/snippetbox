package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	fileServerRoutes "github.com/AlejoGarat/snippetbox/internal/fileserver/routes"
	healtcheckHandler "github.com/AlejoGarat/snippetbox/internal/healthcheck/handlers"
	healtcheckRoutes "github.com/AlejoGarat/snippetbox/internal/healthcheck/routes"
	homeHandler "github.com/AlejoGarat/snippetbox/internal/home/handlers"
	homeRepo "github.com/AlejoGarat/snippetbox/internal/home/repository"
	homeRoutes "github.com/AlejoGarat/snippetbox/internal/home/routes"
	homeService "github.com/AlejoGarat/snippetbox/internal/home/service"
	snippetsHandler "github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
	snippetRepo "github.com/AlejoGarat/snippetbox/internal/snippets/repository"
	snippetsRoutes "github.com/AlejoGarat/snippetbox/internal/snippets/routes"
	snippetService "github.com/AlejoGarat/snippetbox/internal/snippets/service"
	userHandler "github.com/AlejoGarat/snippetbox/internal/users/handlers"
	userRepo "github.com/AlejoGarat/snippetbox/internal/users/repository"
	userRoutes "github.com/AlejoGarat/snippetbox/internal/users/routes"
	userService "github.com/AlejoGarat/snippetbox/internal/users/service"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
)

func main() {
	addr := *flag.String("addr", ":4000", "HTTP network address")
	dsn := *flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httphelpers.NotFound(w)
	})

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	fileServerRoutes.MakeRoutes(router, sessionManager)

	snippetRepo := snippetRepo.NewSnippetRepo(db)
	snippetService := snippetService.NewSnippetService(snippetRepo)

	userRepo := userRepo.NewUserRepo(db)
	userService := userService.NewUserService(userRepo)

	homeRepo := homeRepo.NewHomeRepo(db)
	homeService := homeService.NewHomeService(homeRepo)

	snippetsRoutes.MakeRoutes(router, sessionManager, userRepo, snippetsHandler.New(
		errorLog,
		infoLog,
		snippetService,
		sessionManager,
	))

	homeRoutes.MakeRoutes(router, sessionManager, userRepo, homeHandler.New(
		errorLog,
		infoLog,
		homeService,
		sessionManager,
	))

	userRoutes.MakeRoutes(router, sessionManager, userRepo, userHandler.New(
		errorLog,
		infoLog,
		userService,
		sessionManager,
	))

	healtcheckRoutes.MakeRoutes(router, sessionManager, userRepo, healtcheckHandler.New())

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      router,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", addr)
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
