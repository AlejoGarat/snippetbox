package routes

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
)

type UserRepo interface {
	Exists(id int) (bool, error)
}
type Handler interface {
	HomeView() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(router *httprouter.Router, sessionManager *scs.SessionManager, userRepo UserRepo, handler Handler) {
	dynamic := alice.New(sessionManager.LoadAndSave)
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	router.Handler(http.MethodGet, "/", dynamic.Then(standard.Then(middlewares.SecureHeaders(handler.HomeView()))))
}
