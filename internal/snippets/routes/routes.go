package routes

import (
	"net/http"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
)

type Handler interface {
	SnippetView() func(http.ResponseWriter, *http.Request)
	SnippetCreate() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(mux *http.ServeMux, handler Handler) {
	mux.Handle("/snippet/view", middlewares.RecoverPanic(
		middlewares.LogRequest(
			middlewares.SecureHeaders(
				handler.SnippetView()))))
	mux.Handle("/snippet/create", middlewares.RecoverPanic(
		middlewares.LogRequest(
			middlewares.SecureHeaders(
				handler.SnippetCreate()))))
}
