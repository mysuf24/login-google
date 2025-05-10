package routes

import (
	"login-google/controller"
	"login-google/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.GET("/google/login", controller.GoogleLogin)
		auth.GET("/google/callback", controller.GoogleCallback)
	}

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware())
	user.GET("/profile", controller.Profile)
}
