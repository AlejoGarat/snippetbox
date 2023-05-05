package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (s *Handler) HomeView() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

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
	}
}
