package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/rest/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddlewareAllowsAllOrigins(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCorsMiddlewareAllowsSpecificMethods(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "GET, POST, PUT, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCorsMiddlewareAllowsSpecificHeaders(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestCorsMiddlewareHandlesOptionsRequest(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
