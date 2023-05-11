package routes

import (
	"net/http"

	middlewares "github.com/AlejoGarat/snippetbox/pkg/middlewares"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Handler interface {
	ShowSignup() func(w http.ResponseWriter, r *http.Request)
	Signup() func(http.ResponseWriter, *http.Request)
	ShowLogin() func(http.ResponseWriter, *http.Request)
	Login() func(http.ResponseWriter, *http.Request)
	Logout() func(http.ResponseWriter, *http.Request)
}

func MakeRoutes(router *httprouter.Router, sessionManager *scs.SessionManager, handler Handler) {
	dynamic := alice.New(sessionManager.LoadAndSave)
	standard := alice.New(middlewares.LogRequest, middlewares.LogRequest)
	router.Handler(http.MethodGet, "/user/signup", dynamic.Then(standard.Then(middlewares.SecureHeaders(handler.ShowSignup()))))
	router.Handler(http.MethodPost, "/user/signup", dynamic.Then(standard.Then(middlewares.SecureHeaders(handler.Signup()))))
	router.Handler(http.MethodGet, "/user/login", dynamic.Then(standard.Then(middlewares.SecureHeaders(handler.ShowLogin()))))
	router.Handler(http.MethodPost, "/user/login", dynamic.Then(standard.Then(middlewares.SecureHeaders(handler.Login()))))
	router.Handler(http.MethodPost, "/user/logout", dynamic.Then(standard.Then(middlewares.SecureHeaders(handler.Logout()))))
}
