package repo

import (
	"database/sql"

	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
)

type SnippetRepo struct {
	DB *sql.DB
}

func (sr *SnippetRepo) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := sr.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sr *SnippetRepo) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (sr *SnippetRepo) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
