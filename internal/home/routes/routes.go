package routes

import (
	"net/http"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
)

type Handler interface {
	HomeView() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(mux *http.ServeMux, handler Handler) {
	mux.Handle("/", middlewares.LogRequest(middlewares.SecureHeaders(handler.HomeView())))
}
