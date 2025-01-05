package routes

import (
	"field_archive/server/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDefineRoutes(t *testing.T) {
	// Simulate a test HTTP request
	router := gin.Default()

	h := handlers.RecordingHandler{}
	DefineRoutes(router, &h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "test"}`, w.Body.String())
}
