package repo

import (
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
)

type homeService struct {
	repo HomeRepo
}

type HomeRepo interface {
	Latest() ([]*models.Snippet, error)
}

func NewHomeService(repo HomeRepo) *homeService {
	return &homeService{
		repo: repo,
	}
}

func (hs *homeService) Latest() ([]*models.Snippet, error) {
	return hs.repo.Latest()
}
