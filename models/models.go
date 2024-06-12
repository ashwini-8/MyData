package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Data represents a vehicle 
type Data struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MakeID   int                `json:"Make_ID"`
	MakeName string             `json:"Make_Name"`
}

// Response represents the API response structure
type Response struct {
	Count   int    `json:"Count"`
	Results []Data `json:"Results"`
}