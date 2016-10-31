package models

import "time"
import "gopkg.in/mgo.v2/bson"

//Stroke is an event point on the system
type Stroke struct {
	Location  []float64 `bson:"location" json:"location"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Info      bson.M    `bson:"info" json:"info"`
	CreatedAt time.Time `bson:"created_at" json:"created_time"`
}

type StrokeNear struct {
	TimeRange   int                    `json:"time_range"`
	MaxDistance int                    `json:"max_distance"`
	Stroke      Stroke                 `json:"stroke"`
	Persist     bool                   `json:"persist"`
	Params      map[string]interface{} `json:"params"`
}
