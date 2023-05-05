package repo

import (
	"database/sql"

	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
)

type SnippetRepo struct {
	DB *sql.DB
}

func (sr *SnippetRepo) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (sr *SnippetRepo) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (sr *SnippetRepo) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
