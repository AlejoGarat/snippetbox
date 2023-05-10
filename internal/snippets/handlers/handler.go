package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	"github.com/AlejoGarat/snippetbox/internal/repositoryerrors"
	"github.com/AlejoGarat/snippetbox/internal/serviceerrors"
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
	"github.com/AlejoGarat/snippetbox/pkg/errorhelpers"
	formhelpers "github.com/AlejoGarat/snippetbox/pkg/form"
	"github.com/AlejoGarat/snippetbox/pkg/validator"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	service        SnippetService
	sessionManager *scs.SessionManager
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors errorhelpers.MyErrorMap
	validator.Validator
}

type SnippetService interface {
	Insert(title string, content string, expires int) (int, errorhelpers.MyErrorMap)
	Get(id int) (*models.Snippet, error)
}

func New(errorLog *log.Logger, infoLog *log.Logger, service SnippetService, sessionManager *scs.SessionManager) *handler {
	return &handler{
		errorLog:       errorLog,
		infoLog:        infoLog,
		service:        service,
		sessionManager: sessionManager,
	}
}

func (h *handler) SnippetView() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())

		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil || id < 1 {
			httphelpers.NotFound(w)
			return
		}

		snippet, err := h.service.Get(id)
		if err != nil {
			if errors.Is(err, repositoryerrors.ErrNoRecord) {
				httphelpers.NotFound(w)
			} else {
				httphelpers.ServerError(w, err)
			}
			return
		}

		data := httphelpers.NewTemplateData(r, h.sessionManager)
		data.Snippet = snippet

		httphelpers.Render(w, http.StatusOK, "view.tmpl", templateCache, data)
	}
}

func (h *handler) SnippetCreate() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	var form snippetCreateForm

	return func(w http.ResponseWriter, r *http.Request) {
		err = formhelpers.DecodePostForm(r, &form, formDecoder)
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		form := snippetCreateForm{
			Title:   r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			Expires: expires,
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		id, errMap := h.service.Insert(form.Title, form.Content, form.Expires)
		errs := errMap.Unwrap()

		if len(errs) != 0 {
			err := errors.Join(errs...)
			switch {
			case errors.Is(err, serviceerrors.ErrBlankField),
				errors.Is(err, serviceerrors.ErrExpiresField),
				errors.Is(err, serviceerrors.ErrLongField):

				if !form.Valid() {
					data := httphelpers.NewTemplateData(r, h.sessionManager)
					form.FieldErrors = errMap
					data.Form = form
					httphelpers.Render(w, http.StatusUnprocessableEntity, "create.tmpl", templateCache, data)
					return
				}

			default:
				httphelpers.ServerError(w, err)
			}

			return
		}
		h.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

func (h *handler) SnippetCreateGet() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := httphelpers.NewTemplateData(r, h.sessionManager)

		data.Form = snippetCreateForm{
			Expires: 365,
		}

		httphelpers.Render(w, http.StatusOK, "create.tmpl", templateCache, data)
	}
}
