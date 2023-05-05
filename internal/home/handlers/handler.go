package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
)

type HomeRepo interface {
	Latest() ([]*models.Snippet, error)
}
type handler struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	repo     HomeRepo
}

func New(errorLog *log.Logger, infoLog *log.Logger, repo HomeRepo) *handler {
	return &handler{
		errorLog: errorLog,
		infoLog:  infoLog,
		repo:     repo,
	}
}

func (s *handler) HomeView() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			httphelpers.NotFound(w)
			return
		}

		snippets, err := s.repo.Latest()
		if err != nil {
			httphelpers.ServerError(w, err)
			return
		}

		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v\n", snippet)
		}

		/*

			files := []string{
				"./ui/html/base.tmpl",
				"./ui/html/pages/home.tmpl",
				"./ui/html/partials/nav.tmpl",
			}

			ts, err := template.ParseFiles(files...)
			if err != nil {
				s.ErrorLog.Print(err.Error())
				http.Error(w, "Internal Server Error", 500)
				return
			}

			err = ts.ExecuteTemplate(w, "base", nil)
			if err != nil {
				s.ErrorLog.Print(err.Error())
				http.Error(w, "Internal Server Error", 500)
			}
		*/
	}
}
