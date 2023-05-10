package handlers

import (
	"log"
	"net/http"

	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
	"github.com/alexedwards/scs/v2"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
)

type HomeService interface {
	Latest() ([]*models.Snippet, error)
}
type handler struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	service        HomeService
	sessionManager *scs.SessionManager
}

func New(errorLog *log.Logger, infoLog *log.Logger, service HomeService, sessionManager *scs.SessionManager) *handler {
	return &handler{
		errorLog:       errorLog,
		infoLog:        infoLog,
		service:        service,
		sessionManager: sessionManager,
	}
}

func (s *handler) HomeView() func(w http.ResponseWriter, r *http.Request) {
	templateCache, err := commonmodels.NewTemplateCache()
	if err != nil {
		s.errorLog.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := s.service.Latest()
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}
		data := httphelpers.NewTemplateData(r, s.sessionManager)
		data.Snippets = snippets

		httphelpers.Render(w, http.StatusOK, "home.tmpl", templateCache, data)
	}
}
