package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sourabhch1331/go-jwt-project/controllers"
)

func AuthRouter(incomingRouter *gin.Engine) {
	incomingRouter.POST("/user/sigup", controller.Signup)
	incomingRouter.POST("/user/login", controller.Login)
}
