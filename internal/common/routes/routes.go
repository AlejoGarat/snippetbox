package routes

import "net/http"

type Handler interface {
	HomeView() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(mux *http.ServeMux, handler Handler) {
	mux.HandleFunc("/", handler.HomeView())
}
