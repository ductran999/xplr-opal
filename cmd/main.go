package main

import (
	"net/http"
	"xplr-opal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	authorized := r.Group("/")
	authorized.Use(middlewares.JWTAuthMiddleware(
		"http://localhost:8080/realms/myrealm/protocol/openid-connect/certs",
		"http://localhost:8080/realms/myrealm",
		"account"))

	authorized.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id":          uuid.New(),
			"description": "homework",
		})
	})

	r.Run(":10011")
}
