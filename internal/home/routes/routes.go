package routes

import (
	"net/http"

	"github.com/justinas/alice"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
)

type Handler interface {
	HomeView() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(mux *http.ServeMux, handler Handler) {
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	mux.Handle("/", standard.Then(middlewares.SecureHeaders(handler.HomeView())))
}
