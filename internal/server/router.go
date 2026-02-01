package server

import (
	middleware "github.com/michaelassa01/gomacbot/internal/server/middleware"
	u "github.com/michaelassa01/gomacbot/utils"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRouter(config u.Config) {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.RateLimitMiddleware())

	// Group that ONLY requires API Key
	apiKeyOnly := router.Group("/")
	apiKeyOnly.Use(middleware.APIKeyAuth(config.ApiAccessKey))

	apiKeyOnly.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	s.router = router
}
