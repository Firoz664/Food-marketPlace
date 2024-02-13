package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Food represents a food item in the database.
type Food struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	FoodName  string             `bson:"foodName" json:"foodName" validate:"required,min=2,max=100"`
	Price     float64            `bson:"price" json:"price" validate:"required"`
	FoodImage string             `bson:"foodImage" json:"foodImage" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	FoodID    string             `bson:"foodId" json:"foodId"`
	MenuID    string             `bson:"menuId" json:"menuId" validate:"required"`
}
