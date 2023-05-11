package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	models "github.com/AlejoGarat/snippetbox/internal/models"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
	formhelpers "github.com/AlejoGarat/snippetbox/pkg/form"
	"github.com/AlejoGarat/snippetbox/pkg/validator"
	"github.com/go-playground/form/v4"

	"github.com/alexedwards/scs/v2"
)

type UserService interface {
	ShowSignup() error
	Signup(name string, email string, password string) error
	ShowLogin() error
	Login(email string, password string) (int, error)
	Logout() error
}

type handler struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	service        UserService
	sessionManager *scs.SessionManager
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := httphelpers.NewTemplateData(r, h.sessionManager)
		data.Form = userSignupForm{}

		httphelpers.Render(w, http.StatusOK, "signup.tmpl", templateCache, data)
		fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
	}
}

func (h *handler) Signup() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	var form userSignupForm

	return func(w http.ResponseWriter, r *http.Request) {
		err = formhelpers.DecodePostForm(r, &form, formDecoder)
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

		if !form.Valid() {
			data := httphelpers.NewTemplateData(r, h.sessionManager)
			data.Form = form
			httphelpers.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", templateCache, data)
			return
		}

		err := h.service.Signup(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "Email address is already in use")

				data := httphelpers.NewTemplateData(r, h.sessionManager)
				data.Form = form
				httphelpers.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", templateCache, data)
			} else {
				httphelpers.ServerError(w, err)
			}

			return
		}

		h.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func (h *handler) ShowLogin() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := httphelpers.NewTemplateData(r, h.sessionManager)
		data.Form = userSignupForm{}

		httphelpers.Render(w, http.StatusOK, "login.tmpl", templateCache, data)
	}
}

func (h *handler) Login() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	var form userSignupForm

	return func(w http.ResponseWriter, r *http.Request) {
		err = formhelpers.DecodePostForm(r, &form, formDecoder)
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

		if !form.Valid() {
			data := httphelpers.NewTemplateData(r, h.sessionManager)
			data.Form = form
			httphelpers.Render(w, http.StatusUnprocessableEntity, "login.tmpl", templateCache, data)
			return
		}

		id, err := h.service.Login(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.AddNonFieldError("Email or password is incorrect")

				data := httphelpers.NewTemplateData(r, h.sessionManager)
				data.Form = form
				httphelpers.Render(w, http.StatusUnprocessableEntity, "login.tmpl", templateCache, data)
			} else {
				httphelpers.ServerError(w, err)
			}

			return
		}

		err = h.sessionManager.RenewToken(r.Context())

		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		h.sessionManager.Put(r.Context(), "authenticatedUserID", id)

		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func (h *handler) Logout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.sessionManager.RenewToken(r.Context())
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		h.sessionManager.Remove(r.Context(), "authenticatedUserID")

		h.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
