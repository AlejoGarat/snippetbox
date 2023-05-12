package handlers_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlejoGarat/snippetbox/internal/snippets/handlers"
	"github.com/AlejoGarat/snippetbox/internal/snippets/models"
	"github.com/AlejoGarat/snippetbox/pkg/errorhelpers"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	snippetRoutes "github.com/AlejoGarat/snippetbox/internal/snippets/routes"
)

type UserRepo interface {
	Exists(id int) (bool, error)
}

type ServiceMock struct {
	GetFn    func(id int) (*models.Snippet, error)
	InsertFn func(title string, content string, expires int) (int, errorhelpers.MyErrorMap)
}

func (s *ServiceMock) Get(id int) (*models.Snippet, error) {
	return s.GetFn(id)
}

func (s *ServiceMock) Insert(title string, content string, expires int) (int, errorhelpers.MyErrorMap) {
	return s.InsertFn(title, content, expires)
}

var _ handlers.SnippetService = &ServiceMock{}

func TestSnippetView(t *testing.T) {
	mux := httprouter.New()
	sessionManager := scs.New()
	var userRepo UserRepo

	ts := httptest.NewServer(mux)
	defer ts.Close()

	mock := &ServiceMock{}
	var errorLog *log.Logger
	var infoLog *log.Logger

	handler := handlers.New(errorLog, infoLog, mock, &scs.SessionManager{})

	snippetRoutes.MakeRoutes(mux, sessionManager, userRepo, handler)

	mock.GetFn = func(id int) (*models.Snippet, error) {
		return &models.Snippet{ID: 2}, nil
	}

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := ts.Client().Get(fmt.Sprintf("%s%s", ts.URL, tt.urlPath))
			assert.NoError(t, err, "Failed to send GET request")

			assert.Equal(t, tt.wantCode, response.StatusCode, "Wrong status code")

			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, response.Body, "Wrong body")
			}
		})
	}
}
