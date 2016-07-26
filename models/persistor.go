package models

import (
	"time"

	"github.com/mrkaspa/geoserver/utils"
	"gopkg.in/mgo.v2/bson"
)

type Persistance interface {
	PersistAndFind() chan *StrokeNear
	Persist() chan *Stroke
	UsersFound() chan []Stroke
	Saved() chan bool
}

// Persistor takes care of persisting strokes and finding near ones
type Persistor struct {
	persistAndFind chan *StrokeNear
	persist        chan *Stroke
	usersFound     chan []Stroke
	saved          chan bool
}

func NewPersistor() Persistor {
	persistor := Persistor{
		persistAndFind: make(chan *StrokeNear, 1),
		persist:        make(chan *Stroke, 1),
		usersFound:     make(chan []Stroke, 1),
		saved:          make(chan bool, 1),
	}
	go persistor.run()
	return persistor
}

func (p Persistor) PersistAndFind() chan *StrokeNear {
	return p.persistAndFind
}

func (p Persistor) Persist() chan *Stroke {
	return p.persist
}

func (p Persistor) UsersFound() chan []Stroke {
	return p.usersFound
}

func (p Persistor) Saved() chan bool {
	return p.saved
}

func (p *Persistor) run() {
	defer func() {
		if r := recover(); r != nil {
			utils.Log.Infof("Recovered in persistor.Run()")
		}
	}()
	for {
		select {
		case stroke := <-p.persist:
			utils.Log.Infof("Persistor executing Persist: %s", stroke.UserID)
			p.save(stroke)
		case strokeNear := <-p.persistAndFind:
			utils.Log.Infof("Persistor executing PersistAndFind: %s", strokeNear)
			if p.save(&strokeNear.Stroke) {
				nearUsers, err := p.findNear(strokeNear)
				if err != nil {
					panic(err)
				}
				p.usersFound <- nearUsers
			}
		}
	}
}

func (p *Persistor) save(stroke *Stroke) bool {
	stroke.CreatedAt = time.Now()
	utils.Log.Infof("Persisting %s stroke: %v", stroke.UserID, stroke)
	err := StrokesCollection.Insert(stroke)
	if err != nil {
		p.saved <- false
		return false
	}
	p.saved <- true
	return true
}

func (p *Persistor) findNear(strokeNear *StrokeNear) ([]Stroke, error) {
	results := []Stroke{}
	query := buildQuery(strokeNear)
	err := StrokesCollection.Find(query).All(&results)
	utils.Log.Infof("Query executed by %s: %v", strokeNear.Stroke.UserID, query)
	utils.Log.Infof("Actor %s found matches %d", strokeNear.Stroke.UserID, len(results))
	return results, err
}

func buildQuery(strokeNear *StrokeNear) bson.M {
	query := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": strokeNear.Stroke.Location,
				},
				"$maxDistance": strokeNear.MaxDistance,
			},
		},
		"user_id": bson.M{
			"$ne": strokeNear.Stroke.UserID,
		},
	}

	if strokeNear.TimeRange > 0 {
		query["created_at"] = bson.M{
			"$gte": time.Now().Add(-1 * time.Duration(strokeNear.TimeRange) * time.Second),
			"$lte": time.Now(),
		}
	}

	return query
}
