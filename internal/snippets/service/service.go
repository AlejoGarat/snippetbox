package repo

import (
	"strings"
	"unicode/utf8"

	"github.com/AlejoGarat/snippetbox/internal/serviceerrors"
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
	"github.com/AlejoGarat/snippetbox/pkg/errorhelpers"
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

func (ss *snippetService) Insert(title string, content string, expires int) (int, errorhelpers.MyErrorMap) {
	var fieldErrors errorhelpers.MyErrorMap

	if strings.TrimSpace(title) == "" {
		fieldErrors = errorhelpers.JoinErrorMap(fieldErrors, "title", serviceerrors.ErrBlankField)
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors = errorhelpers.JoinErrorMap(fieldErrors, "title", serviceerrors.ErrLongField)
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors = errorhelpers.JoinErrorMap(fieldErrors, "content", serviceerrors.ErrBlankField)
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors = errorhelpers.JoinErrorMap(fieldErrors, "expires", serviceerrors.ErrExpiresField)
	}

	if len(fieldErrors) > 0 {
		return invalidId, fieldErrors
	}

	id, err := ss.repo.Insert(title, content, expires)

	errorhelpers.JoinErrorMap(fieldErrors, "insert", err)
	return id, fieldErrors
}

func (ss *snippetService) Get(id int) (*models.Snippet, error) {
	return ss.repo.Get(id)
}
