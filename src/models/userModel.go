package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProfileType represents the possible profile types
type ProfileType string

const (
	UserType       ProfileType = "USER"
	RiderType      ProfileType = "RIDER"
	RestaurantType ProfileType = "RESTAURANT"
)

func ValidateProfileType(profileType ProfileType) error {
	switch profileType {
	case UserType, RiderType, RestaurantType:
		return nil
	default:
		return errors.New("invalid profile type")
	}
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName    *string            `json:"firstName,omitempty" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastName,omitempty" validate:"required,min=2,max=100"`
	Password     *string            `json:"password,omitempty" validate:"required,min=6"`
	IsAdmin      bool               `json:"isAdmin,omitempty"`
	ProfileType  *ProfileType       `json:"profileType,omitempty" bson:"profileType,omitempty" validate:"required"`
	Email        *string            `json:"email,omitempty" validate:"email,required"`
	ProfileImage *string            `json:"profileImage,omitempty" bson:"profileImage,omitempty"`
	MobileNumber *string            `json:"mobileNumber,omitempty" validate:"required" bson:"mobileNumber,omitempty"`
	Token        *string            `json:"token,omitempty"`
	RefreshToken *string            `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	UserID       string             `json:"userId,omitempty" bson:"userId,omitempty"`
	IsVerify     bool               `json:"isVerify,omitempty" bson:"isVerify,omitempty"`
	IsComplete   bool               `json:"isComplete,omitempty" bson:"isComplete,omitempty"`
	LastActive   time.Time          `json:"lastActive,omitempty" bson:"lastActive,omitempty"`
}
