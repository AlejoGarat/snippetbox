package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
)

type Handler interface {
	HomeView() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(router *httprouter.Router, handler Handler) {
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	router.Handler(http.MethodGet, "/", standard.Then(middlewares.SecureHeaders(handler.HomeView())))
}
