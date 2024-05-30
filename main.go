package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	routes "github.com/sourabhch1331/go-jwt-project/routes"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRouter(router)
	routes.AuthRouter(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statusCode": 200,
			"message":    "acces granted for api-1",
			"data":       nil,
		})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statusCode": 200,
			"message":    "acces granted for api-2",
			"data":       nil,
		})
	})

	router.Run(":" + port)
}
