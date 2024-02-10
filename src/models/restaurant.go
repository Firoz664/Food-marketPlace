package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Location    Location           `bson:"location" json:"location"`
	Menu        []Menu             `bson:"menu" json:"menu"`
}

// Location represents the schema for a restaurant's location
type Location struct {
	Address     string  `bson:"address" json:"address"`
	City        string  `bson:"city" json:"city"`
	State       string  `bson:"state" json:"state"`
	ZipCode     string  `bson:"zipcode" json:"zipcode"`
	Coordinates GeoJSON `bson:"coordinates" json:"coordinates"`
}

// GeoJSON represents geographical coordinates
type GeoJSON struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
