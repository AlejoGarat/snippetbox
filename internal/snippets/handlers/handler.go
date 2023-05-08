package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	"github.com/AlejoGarat/snippetbox/internal/repositoryerrors"
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
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
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		data := &commonmodels.TemplateData{
			Snippet: snippet,
		}

		templateCache, err := commonmodels.NewTemplateCache()
		if err != nil {
			h.errorLog.Fatal(err)
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			httphelpers.ServerError(w, err)
		}

		httphelpers.Render(w, http.StatusOK, "view.tmpl", templateCache, &commonmodels.TemplateData{
			Snippet: snippet,
		})
	}
}

func (h *handler) SnippetCreate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		//Simulate incoming data
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

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
