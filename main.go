package main

import (
	"log"
	"login-google/config"
	"login-google/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitGoogleOAuth()

	r := gin.Default()

	routes.RegisterAuthRoutes(r)

	log.Println("Server running at http://localhost:8080")

	r.Run(":8080")

}
