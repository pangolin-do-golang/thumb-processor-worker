package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/rest/handler"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerHandlerReturnsNotFound(t *testing.T) {
	router := gin.Default()
	handler.RegisterSwaggerHandlers(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/swagger/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
