package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	healtcheckRoutes "github.com/AlejoGarat/snippetbox/internal/healthcheck/routes"
)

type UserRepo interface {
	Exists(id int) (bool, error)
}

func TestPing(t *testing.T) {
	mux := httprouter.New()
	sessionManager := scs.New()
	var userRepo UserRepo

	healtcheckRoutes.MakeRoutes(mux, sessionManager, userRepo, New())

	ts := httptest.NewServer(mux)
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, "OK", string(body))
}
