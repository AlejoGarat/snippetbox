package routes

import (
	"net/http"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Handler interface {
	SnippetView() func(http.ResponseWriter, *http.Request)
	SnippetCreate() func(http.ResponseWriter, *http.Request)
	SnippetCreateGet() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(router *httprouter.Router, handler Handler) {
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	router.Handler(http.MethodGet, "/snippet/view/:id", standard.Then(middlewares.SecureHeaders(handler.SnippetView())))
	router.Handler(http.MethodPost, "/snippet/create", standard.Then(middlewares.SecureHeaders(handler.SnippetCreate())))
	router.Handler(http.MethodGet, "/snippet/create", standard.Then(middlewares.SecureHeaders(handler.SnippetCreateGet())))
}
