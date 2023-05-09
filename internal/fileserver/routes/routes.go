package routes

import (
	"net/http"

	"github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func MakeRoutes(router *httprouter.Router) {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	router.Handler(http.MethodGet, "/static/", standard.Then(middlewares.SecureHeaders(http.StripPrefix("/static", fileServer).ServeHTTP)))
}
