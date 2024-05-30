package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name" validate:"required min=2 max=100"`
	Password     *string            `json:"password" validate:"required min=6"`
	Email        *string            `json:"email" validate:"required email"`
	Phone        *string            `json:"phone" validate:"required min=10 max=10"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken"`
	UserType     *string            `json:"userType" validate:"required eq=ADMIN|eq=USER"`
	CreatedAt    time.Time          `json:"createdAt"`
	UserId       *string            `json:"userId"`
}
