package handlers

import (
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
		snippets, err := s.service.Latest()
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		data := httphelpers.NewTemplateData(r)
		data.Snippets = snippets

		templateCache, err := commonmodels.NewTemplateCache()
		if err != nil {
			s.errorLog.Fatal(err)
		}

		httphelpers.Render(w, http.StatusOK, "home.tmpl", templateCache, data)
	}
}
