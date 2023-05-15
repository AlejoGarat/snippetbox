package routes

import (
	"net/http"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type UserRepo interface {
	Exists(id int) (bool, error)
}

type Handler interface {
	Ping() func(w http.ResponseWriter, r *http.Request)
}

func MakeRoutes(router *httprouter.Router, sessionManager *scs.SessionManager, userRepo UserRepo, handler Handler) {
	mids := alice.New(sessionManager.LoadAndSave, middlewares.NoSurf, middlewares.LogRequest, middlewares.LogRequest,
		middlewares.RecoverPanic)

	router.Handler(http.MethodGet, "/ping",
		mids.Then(
			middlewares.SecureHeaders(
				handler.Ping(),
			),
		),
	)
}
