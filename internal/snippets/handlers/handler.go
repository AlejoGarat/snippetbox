package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	repo "github.com/AlejoGarat/snippetbox/internal/snippets/repository"
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

		w.Write([]byte("Create a new snippet..."))
	}
}
