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
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	service  SnippetService
}

type SnippetService interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*models.Snippet, error)
}

func New(errorLog *log.Logger, infoLog *log.Logger, service SnippetService) *handler {
	return &handler{
		errorLog: errorLog,
		infoLog:  infoLog,
		service:  service,
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

		data := httphelpers.NewTemplateData(r)
		data.Snippet = snippet

		httphelpers.Render(w, http.StatusOK, "view.tmpl", templateCache, data)
	}
}

func (h *handler) SnippetCreate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			httphelpers.ClientError(w, http.StatusBadRequest)
			return
		}

		id, err := h.service.Insert(title, content, expires)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrBlankField),
				errors.Is(err, serviceerrors.ErrExpiresField),
				errors.Is(err, serviceerrors.ErrLongField):
				httphelpers.BadRequestError(w, err)
			default:
				httphelpers.ServerError(w, err)
			}

			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

func (h *handler) SnippetCreateGet() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		h.errorLog.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := httphelpers.NewTemplateData(r)

		httphelpers.Render(w, http.StatusOK, "create.tmpl", templateCache, data)
	}
}
