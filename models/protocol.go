package models

type PersistWithResponse struct {
	Stroke   Stroke
	Response chan bool
}

type PersistAndFindWithResponse struct {
	StrokeNear    StrokeNear
	SavedResponse chan bool
	UsersResponse chan []Stroke
}
