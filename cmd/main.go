package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "Hello there, server is running",
		})
	})
	router.Run() // Listening and serving on 0.0.0.0:8080
}
