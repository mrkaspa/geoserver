package models

import "gopkg.in/mgo.v2/bson"

type mockPersistor struct {
}

func NewMockPersistor() Persistance {
	return mockPersistor{}
}

func (p mockPersistor) PersistAndFind(sn StrokeNear) ([]Stroke, error) {
	return []Stroke{
		{
			UserID:   "a2",
			Info:     bson.M{},
			Location: []float64{-79.38066843, 43.65483486},
		},
	}, nil
}

func (p mockPersistor) Persist(s Stroke) error {
	return nil
}

func (p mockPersistor) FindStrokes(username string) ([]Stroke, error) {
	return []Stroke{
		{
			UserID:   "a1",
			Info:     bson.M{},
			Location: []float64{-79.38066843, 43.65483486},
		},
	}, nil
}
