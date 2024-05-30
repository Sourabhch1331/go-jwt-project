package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HashPassword() {

}

func VerifyPassword() {

}

func Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "user logged in!",
		"data":       nil,
	})
}

func Signup(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "user signed up!",
		"data":       nil,
	})
}

func GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "All users",
		"data":       nil,
	})
}

func GetUser(ctx *gin.Context) {
	userId, _ := ctx.Params.Get("userId")

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "user with id " + userId,
		"data":       nil,
	})
}
