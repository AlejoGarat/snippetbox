package routes

import (
	"net/http"

	"github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/justinas/alice"
)

func MakeRoutes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	mux.Handle("/static/", standard.Then(middlewares.SecureHeaders(http.StripPrefix("/static", fileServer).ServeHTTP)))
}
