package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(ctx *gin.Context, role string) error {
	// some other middleware will put this value to context about current user
	userType := ctx.GetString("userType")

	if userType != role {
		return errors.New("unauthorized access to resource")
	}

	return nil
}

func MatchUserTypeToUId(ctx *gin.Context, userId string) error {
	// some other middleware will put this value to context about current user
	userType := ctx.GetString("userType")
	currUserId := ctx.GetString("uId")

	if userType == "USER" && userId != currUserId {
		return errors.New("unauthorized access to resource")
	}

	return nil
}
