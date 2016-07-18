package models

import "time"

//Stroke is an event point on the system
type Stroke struct {
	Location  []float64 `bson:"location"`
	UserID    string    `bson:"user_id"`
	Info      string    `bson:"info"`
	CreatedAt time.Time `bson:"created_at"`
}
