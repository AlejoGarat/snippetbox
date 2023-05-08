package httphelpers

import (
	"fmt"
	"net/http"
	"text/template"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
)

func Render(w http.ResponseWriter, status int, page string, templateCache map[string]*template.Template, data *commonmodels.TemplateData) {
	ts, ok := templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		ServerError(w, err)
	}
}
