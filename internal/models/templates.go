package models

import "github.com/AlejoGarat/snippetbox/internal/snippets/models"

type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
