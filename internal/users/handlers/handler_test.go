package handlers_test

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/AlejoGarat/snippetbox/internal/users/handlers"
	userRoutes "github.com/AlejoGarat/snippetbox/internal/users/routes"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const (
	validName     = "Bob"
	validPassword = "validPa$$word"
	validEmail    = "bob@example.com"
	formTag       = "<form action='/user/signup' method='POST' novalidate>"
)

type UserRepo interface {
	Exists(id int) (bool, error)
}

type ServiceMock struct {
	SignupFn     func(name string, email string, password string) error
	ShowLoginFn  func() error
	LoginFn      func(email string, password string) (int, error)
	LogoutFn     func() error
	ShowSignupFn func() error
}

func (s *ServiceMock) ShowSignup() error {
	return nil
}

func (s *ServiceMock) Signup(name string, email string, password string) error {
	return s.SignupFn(name, email, password)
}

func (s *ServiceMock) ShowLogin() error {
	return nil
}

func (s *ServiceMock) Login(email string, password string) (int, error) {
	return s.LoginFn(email, password)
}

func (s *ServiceMock) Logout() error {
	return nil
}

var _ handlers.UserService = &ServiceMock{}

func TestUserSignup(t *testing.T) {
	mux := httprouter.New()
	sessionManager := scs.New()
	var userRepo UserRepo

	ts := httptest.NewServer(mux)
	defer ts.Close()

	mock := &ServiceMock{}
	var errorLog *log.Logger
	var infoLog *log.Logger

	handler := handlers.New(errorLog, infoLog, mock, &scs.SessionManager{})

	userRoutes.MakeRoutes(mux, sessionManager, userRepo, handler)

	validCSRFToken := extractCSRFToken(t, "")

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			response, err := ts.Client().Get(fmt.Sprintf("%s%s", ts.URL, "/user/signup"))
			assert.NoError(t, err, "Failed to send GET request")

			assert.Equal(t, tt.wantCode, response.StatusCode)

			if tt.wantFormTag != "" {
				assert.Equal(t, response.Body, tt.wantFormTag)
			}
		})
	}
}

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}
