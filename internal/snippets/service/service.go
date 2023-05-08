package repo

import (
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
)

type snippetService struct {
	repo SnippetRepo
}

type SnippetRepo interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*models.Snippet, error)
}

func NewSnippetService(repo SnippetRepo) *snippetService {
	return &snippetService{
		repo: repo,
	}
}

func (ss *snippetService) Insert(title string, content string, expires int) (int, error) {
	return ss.repo.Insert(title, content, expires)
}

func (ss *snippetService) Get(id int) (*models.Snippet, error) {
	return ss.repo.Get(id)
}
