package tests

import (
	"hashtechy/src/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerInitialization(t *testing.T) {
	mux := server.Server()
	assert.NotNil(t, mux, "Server should be initialized")
}

func TestSwaggerHandler(t *testing.T) {
	mux := server.Server()
	server.AddSwaggerHandler(mux)

	req, err := http.NewRequest("GET", "/swagger/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Swagger handler should return 200 OK")
}
