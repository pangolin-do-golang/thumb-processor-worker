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
