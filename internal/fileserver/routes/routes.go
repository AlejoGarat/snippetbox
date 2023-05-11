package routes

import (
	"net/http"

	"github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func MakeRoutes(router *httprouter.Router, sessionManager *scs.SessionManager) {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mids := alice.New(sessionManager.LoadAndSave, middlewares.LogRequest, middlewares.LogRequest)
	router.Handler(http.MethodGet, "/static/",
		mids.Then(
			middlewares.SecureHeaders(
				http.StripPrefix("/static", fileServer).ServeHTTP,
			),
		),
	)
}
