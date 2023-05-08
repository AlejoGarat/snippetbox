package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/AlejoGarat/snippetbox/internal/snippets/models"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
)

type HomeService interface {
	Latest() ([]*models.Snippet, error)
}
type handler struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	service  HomeService
}

func New(errorLog *log.Logger, infoLog *log.Logger, service HomeService) *handler {
	return &handler{
		errorLog: errorLog,
		infoLog:  infoLog,
		service:  service,
	}
}

func (s *handler) HomeView() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			httphelpers.NotFound(w)
			return
		}

		snippets, err := s.service.Latest()
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		data := &commonmodels.TemplateData{
			Snippets: snippets,
		}

		// Pass in the templateData struct when executing the template.
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			httphelpers.ServerError(w, err)
		}
	}
}
