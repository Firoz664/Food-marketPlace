package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Driver represents the schema for drivers collection
type Driver struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name"`
	PhoneNumber string             `json:"phoneNumber"`
	Email       string             `json:"email"`
	License     string             `json:"license"`
	Status      string             `json:"status"`
	VehicleID   primitive.ObjectID `json:"vehicleId,omitempty"`
}

// Vehicle represents the schema for vehicles collection
type Vehicle struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Model       string             `json:"model"`
	Make        string             `json:"make"`
	Year        int                `bjson:"year"`
	PlateNumber string             `json:"plateNumber"`
	Color       string             `json:"color"`
	OwnerID     primitive.ObjectID `bson:"ownerId,omitempty"`
}
