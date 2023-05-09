package repo

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/AlejoGarat/snippetbox/internal/serviceerrors"
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
)

const (
	invalidId = -1
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
	var fieldErrors []error

	if strings.TrimSpace(title) == "" {
		fieldErrors = append(fieldErrors, serviceerrors.ErrBlankField)
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors = append(fieldErrors, serviceerrors.ErrLongField)
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors = append(fieldErrors, serviceerrors.ErrBlankField)
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors = append(fieldErrors, serviceerrors.ErrExpiresField)
	}

	if len(fieldErrors) > 0 {
		return invalidId, errors.Join(fieldErrors...)
	}

	return ss.repo.Insert(title, content, expires)
}

func (ss *snippetService) Get(id int) (*models.Snippet, error) {
	return ss.repo.Get(id)
}
