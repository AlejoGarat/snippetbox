package routes

import (
	"net/http"

	"github.com/AlejoGarat/snippetbox/pkg/middlewares"
)

func MakeRoutes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", middlewares.SecureHeaders(
		http.StripPrefix("/static", fileServer).ServeHTTP),
	)
}
