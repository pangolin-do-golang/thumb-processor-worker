package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerHandlers(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
