package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastName" validate:"required,min=2,max=100"`
	Password     *string            `json:"Password" validate:"required,min=6"`
	IsAdmin      bool               `json:"isAdmin"`
	Email        *string            `json:"email" validate:"email,required"`
	ProfileImage *string            `json:"profileImage"`
	MobileNumber *string            `json:"mobileNumber" validate:"required"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserID       string             `json:"userId"`
	IsVerify     bool               `json:"isVerify"`
	IsComplete   bool               `json:"isComplete"`
	LastActive   time.Time          `json:"lastActive"`
}
