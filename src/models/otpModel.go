package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `json:"email,omitempty" validate:"email"`
	MobileNumber string             `json:"mobileNumber,omitempty"`
	OtpCode      int64              `json:"otpCode"   bson:"otpCode"`
	IsVerify     bool               `json:"isVerify"  bson:"isVerify,default=false"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	ExpiredAt    time.Time          `json:"expiredAt" bson:"expiredAt"`
}
