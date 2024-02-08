package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID         primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"` // Assuming you have a corresponding User struct
	PaymentID      string             `bson:"paymentId,omitempty" json:"paymentId,omitempty"`
	PaymentStatus  string             `bson:"paymentStatus,omitempty" json:"paymentStatus,omitempty"`
	Amount         string             `bson:"amount,omitempty" json:"amount,omitempty"`
	OriginalAmount string             `bson:"originalAmount,omitempty" json:"originalAmount,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
