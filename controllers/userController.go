package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sourabhch1331/go-jwt-project/database"
	helper "github.com/sourabhch1331/go-jwt-project/helpers"
	"github.com/sourabhch1331/go-jwt-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword() {

}

func VerifyPassword(hashedPassword string, planPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(planPassword))
	return err == nil
}

func Login(ctx *gin.Context) {

	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user models.User
	var foundUser models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "fail",
			"statusCode": http.StatusBadRequest,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	if err := userCollection.FindOne(c, bson.M{"email": *user.Email}).Decode(&foundUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "fail",
			"statusCode": http.StatusBadRequest,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	if passwordIsValid := VerifyPassword(*foundUser.Password, *user.Password); !passwordIsValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "fail",
			"statusCode": http.StatusUnauthorized,
			"message":    "wrong password!",
			"data":       nil,
		})
		return
	}

	token, refreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.UserType, *foundUser.UserId)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	foundUser.Token = &token
	foundUser.RefreshToken = &refreshToken

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refreshtoken", Value: refreshToken})

	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = userCollection.UpdateOne(c, bson.M{"email": *foundUser.Email}, bson.D{{Key: "$set", Value: updateObj}}, &opts)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "user logged in!",
		"data":       foundUser,
	})
}

func Signup(ctx *gin.Context) {

	// define context
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	// try to bind req.body to user
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "fail",
			"statusCode": http.StatusBadRequest,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	// validate the user data recived from body
	validationError := validate.Struct(user)

	if validationError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "fail",
			"statusCode": http.StatusBadRequest,
			"message":    validationError.Error(),
			"data":       nil,
		})
		return
	}

	// count if user with email already exist
	count, err := userCollection.CountDocuments(c, bson.M{"email": *user.Email})

	// if err return
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	//  if user with email already exist, return
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "fail",
			"statusCode": http.StatusBadRequest,
			"message":    "User already exist with email: " + *user.Email,
			"data":       nil,
		})
	}

	// count if user with phone already exist
	count, err = userCollection.CountDocuments(c, bson.M{"phone": *user.Phone})

	// if err return
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	//  if user with phone already exist, return
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "fail",
			"statusCode": http.StatusBadRequest,
			"message":    "User already exist with email: " + *user.Email,
			"data":       nil,
		})
	}

	// generate CreatedAt time stamp and userId

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Id = primitive.NewObjectID()
	uId := user.Id.Hex()
	user.UserId = &uId

	// gen token
	token, refershToken, err := helper.GenerateAllTokens(*user.Email, *user.Name, *user.UserType, *user.UserId)

	fmt.Println("Here")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    "failed to generate jwt token",
			"data":       nil,
		})
	}

	user.Token = &token
	user.RefreshToken = &refershToken

	HashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), 14)
	hp := string(HashedPassword)
	user.Password = &hp

	// insert in db
	_, err = userCollection.InsertOne(c, &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    "User not created!",
			"data":       nil,
		})
		return
	}

	// send inserted user
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "user signed up",
		"data":       user,
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
	// get user id from param
	userId := ctx.Param("userId")

	// if user role is not admin and he/she is trying to acces other user
	// then return error
	if err := helper.MatchUserTypeToUId(ctx, userId); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "fail",
			"statusCode": http.StatusUnauthorized,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	var user models.User

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find user by userId and store in user(models.User)
	err := userCollection.FindOne(c, bson.M{"userid": userId}).Decode(&user)

	// if error fetching user return
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "fail",
			"statusCode": http.StatusInternalServerError,
			"message":    err.Error(),
			"data":       nil,
		})
		return
	}

	// everthing went well return the user with userId

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "user with id " + userId,
		"data":       user,
	})
}
