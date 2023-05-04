package routes

import "net/http"

type Handler interface {
	SnippetView() func(http.ResponseWriter, *http.Request)
	SnippetCreate() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(mux *http.ServeMux, handler Handler) {
	mux.HandleFunc("/snippet/view", handler.SnippetView())
	mux.HandleFunc("/snippet/create", handler.SnippetCreate())
}
