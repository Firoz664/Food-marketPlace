package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApprovalStatus represents the status of approval
type ApprovalStatus string

const (
	Approve   ApprovalStatus = "Approve"
	Hold      ApprovalStatus = "Hold"
	Pending   ApprovalStatus = "Pending"
	Rejected  ApprovalStatus = "Rejected"
	Suspended ApprovalStatus = "Suspended"
)

// ValidateApprovalStatus validates the ApprovalStatus
func ValidateApprovalStatus(status ApprovalStatus) error {
	switch status {
	case Hold, Pending, Approve, Suspended, Rejected:
		return nil
	default:
		return errors.New("invalid approval status")
	}
}

// Location represents the schema for a restaurant's location
type Location struct {
	Street          string `bson:"street" json:"street"`
	LocalityVerbose string `bson:"localityVerbose" json:"localityVerbose"`
	City            string `bson:"city" json:"city"`
	State           string `bson:"state" json:"state"`
	Coordinates     struct {
		Latitude  string `bson:"latitude" json:"latitude"`
		Longitude string `bson:"longitude" json:"longitude"`
	} `bson:"coordinates" json:"coordinates"`
}

// Restaurant represents a restaurant entity
type Restaurant struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            string             `bson:"userId,omitempty" json:"userId"`
	RestaurantName    string             `bson:"restaurantName" json:"restaurantName"`
	BannerImage       string             `bson:"bannerImage" json:"bannerImage"`
	Description       string             `bson:"description" json:"description"`
	Status            ApprovalStatus     `bson:"status" json:"status"`
	Location          Location           `bson:"location" json:"location"`
	Cuisines          []string           `bson:"cuisines" json:"cuisines"`
	IsDelivering      bool               `bson:"isDelivering" json:"isDelivering"`
	HasOnlineDelivery bool               `bson:"hasOnlineDelivery" json:"hasOnlineDelivery"`
	Menu              []Menu             `bson:"menu" json:"menu"`
}
