package routes

import (
	"net/http"

	"github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/AlejoGarat/snippetbox/ui"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func MakeRoutes(router *httprouter.Router, sessionManager *scs.SessionManager) {
	fileServer := http.FileServer(http.FS(ui.Files))
	mids := alice.New(sessionManager.LoadAndSave, middlewares.LogRequest, middlewares.LogRequest, middlewares.RecoverPanic)
	router.Handler(http.MethodGet, "/static/*filepath",
		mids.Then(
			middlewares.SecureHeaders(
				http.StripPrefix("/static", fileServer).ServeHTTP,
			),
		),
	)
}
