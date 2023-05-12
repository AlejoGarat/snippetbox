package handlers

import (
	"net/http"
)

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) Ping() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}
