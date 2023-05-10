package httphelpers

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"

	commonmodels "github.com/AlejoGarat/snippetbox/internal/models"
	"github.com/alexedwards/scs/v2"
)

func Render(w http.ResponseWriter, status int, page string, templateCache map[string]*template.Template, data *commonmodels.TemplateData) {
	ts, ok := templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func NewTemplateData(r *http.Request, sessionManager *scs.SessionManager) *commonmodels.TemplateData {
	return &commonmodels.TemplateData{
		CurrentYear: time.Now().Year(),
		Flash:       sessionManager.PopString(r.Context(), "flash"),
	}
}
