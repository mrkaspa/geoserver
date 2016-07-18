package models

import "time"

//Stroke is an event point on the system
type Stroke struct {
	Location  []float64 `bson:"location" json:"location"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Info      string    `bson:"info" json:"info"`
	CreatedAt time.Time `bson:"created_at" json:"created_time"`
}

type StrokeNear struct {
	TimeRange   int    `json:"time_range"`
	MaxDistance int    `json:"max_distance"`
	Stroke      Stroke `json:"stroke"`
}
