package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sourabhch1331/go-jwt-project/controllers"
)

func UserRouter(incomingRouter *gin.Engine) {
	// incomingRouter.Use(middleware.Authenticate())

	incomingRouter.GET("/user", controller.GetUsers)
	incomingRouter.GET("/user/:userId", controller.GetUser)

}
