package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	database "github.com/sourabhch1331/go-jwt-project/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email    string
	Name     string
	UId      string
	UserType string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, name string, userType string, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		UserType: userType,
		UId:      userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
