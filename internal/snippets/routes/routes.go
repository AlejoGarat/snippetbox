package routes

import (
	"net/http"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/justinas/alice"
)

type Handler interface {
	SnippetView() func(http.ResponseWriter, *http.Request)
	SnippetCreate() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(mux *http.ServeMux, handler Handler) {
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	mux.Handle("/snippet/view", standard.Then(middlewares.SecureHeaders(handler.SnippetView())))
	mux.Handle("/snippet/create", standard.Then(middlewares.SecureHeaders(handler.SnippetCreate())))
}
