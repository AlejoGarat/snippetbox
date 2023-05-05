package handlers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	repo "github.com/AlejoGarat/snippetbox/internal/snippets/repository"
	httphelpers "github.com/AlejoGarat/snippetbox/pkg"
)

type Handler struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Repo     *repo.SnippetRepo
}

func (h *Handler) SnippetView() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	}
}

func (h *Handler) SnippetCreate() func(w http.ResponseWriter, r *http.Request) {
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

		id, err := h.Repo.Insert(title, content, expires)
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
			h.ErrorLog.Output(2, trace)
			httphelpers.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
