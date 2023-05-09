package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	"github.com/AlejoGarat/snippetbox/internal/repositoryerrors"
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

		templateCache, err := commonmodels.NewTemplateCache()
		if err != nil {
			h.errorLog.Fatal(err)
		}

		httphelpers.Render(w, http.StatusOK, "view.tmpl", templateCache, data)
	}
}

func (h *handler) SnippetCreate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simulate incoming data
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		id, err := h.service.Insert(title, content, expires)
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
			h.errorLog.Output(2, trace)
			httphelpers.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

func (h *handler) SnippetCreateGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Display the form for creating a new snippet..."))
	}
}
