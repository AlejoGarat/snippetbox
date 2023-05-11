package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type UserService interface {
	ShowSignup() error
	Signup() error
	ShowLogin() error
	Login() error
	Logout() error
}

type handler struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	service        UserService
	sessionManager *scs.SessionManager
}

func New(errorLog *log.Logger, infoLog *log.Logger, service UserService, sessionManager *scs.SessionManager) *handler {
	return &handler{
		errorLog:       errorLog,
		infoLog:        infoLog,
		service:        service,
		sessionManager: sessionManager,
	}
}

func (h *handler) ShowSignup() func(w http.ResponseWriter, r *http.Request) {
	/*templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}*/

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
	}
}

func (h *handler) Signup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create a new user...")
	}
}

func (h *handler) ShowLogin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Display a HTML form for logging in a user...")
	}
}

func (h *handler) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Authenticate and login the user...")
	}
}

func (h *handler) Logout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Logout the user...")
	}
}
