package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/rest/handler"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/rest/middleware"
)

type RestServer struct {
}

type RestServerOptions struct {
}

func NewRestServer(_ *RestServerOptions) *RestServer {
	return &RestServer{}
}

// RegisterHealthCheck registers the health check endpoint with the given Gin engine.
//
// @Tags Health Check
// @Summary      Health Check
// @Description  Checks the health status of the application.
// @Produce      json
// @Success     200 {object} map[string]interface{}
// @Router      /health [get]
func (rs RestServer) Serve() {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	handler.RegisterSwaggerHandlers(r)
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
