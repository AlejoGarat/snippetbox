package repo

import (
	"database/sql"
	"errors"

	"github.com/AlejoGarat/snippetbox/internal/repositoryerrors"
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
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := sr.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositoryerrors.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (sr *SnippetRepo) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
