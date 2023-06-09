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
	ShowSignup() func(w http.ResponseWriter, r *http.Request)
	Signup() func(http.ResponseWriter, *http.Request)
	ShowLogin() func(http.ResponseWriter, *http.Request)
	Login() func(http.ResponseWriter, *http.Request)
	Logout() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(router *httprouter.Router, sessionManager *scs.SessionManager, userRepo UserRepo, handler Handler) {
	mids := alice.New(sessionManager.LoadAndSave, middlewares.NoSurf, middlewares.LogRequest, middlewares.LogRequest,
		middlewares.RecoverPanic)

	router.Handler(http.MethodGet, "/user/signup",
		mids.Then(
			middlewares.SecureHeaders(
				handler.ShowSignup(),
			),
		),
	)

	router.Handler(http.MethodPost, "/user/signup",
		mids.Then(
			middlewares.SecureHeaders(
				handler.Signup(),
			),
		),
	)

	router.Handler(http.MethodGet, "/user/login",
		mids.Then(
			middlewares.SecureHeaders(
				handler.ShowLogin(),
			),
		),
	)

	router.Handler(http.MethodPost, "/user/login",
		mids.Then(
			middlewares.SecureHeaders(
				handler.Login(),
			),
		),
	)

	router.Handler(http.MethodPost, "/user/logout",
		mids.Then(
			middlewares.SecureHeaders(
				handler.Logout(),
			),
		),
	)
}
