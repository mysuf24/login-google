package routes

import (
	"login-google/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.GET("/google/login", controller.GoogleLogin)
		auth.GET("/google/callback", controller.GoogleCallback)
	}
}
