package models

import (
	"time"

	"github.com/mrkaspa/geoserver/utils"
	"gopkg.in/mgo.v2/bson"
)

const (
	// max distance for the query in meters
	maxDistance = 5
	// seconds range to do match
	secondsRange = 3
)

type Persistor struct {
	PersistAndFind chan *Stroke
	Persist        chan *Stroke
	UsersFound     chan []Stroke
	Saved          chan bool
}

func NewPersistor() *Persistor {
	persistor := &Persistor{
		PersistAndFind: make(chan *Stroke, 1),
		Persist:        make(chan *Stroke, 1),
		UsersFound:     make(chan []Stroke, 1),
		Saved:          make(chan bool, 1),
	}
	go persistor.run()
	return persistor
}

func (p *Persistor) run() {
	defer close(p.Persist)
	defer close(p.PersistAndFind)

	select {
	case stroke := <-p.Persist:
		utils.Log.Infof("Persistor executing Persist: %s", stroke.UserID)
		if err := p.save(stroke); err != nil {
			p.Saved <- false
			return
		}
		p.Saved <- true
	case stroke := <-p.PersistAndFind:
		utils.Log.Infof("Persistor executing PersistAndFind: %s", stroke.UserID)
		if err := p.save(stroke); err != nil {
			p.Saved <- false
			return
		}
		p.Saved <- true
		nearUsers, err := p.findNear(stroke)
		if err != nil {
			panic(err)
		}
		p.UsersFound <- nearUsers
	}
}

func (p *Persistor) save(stroke *Stroke) error {
	stroke.CreatedAt = time.Now()
	utils.Log.Infof("Persisting %s stroke: %v", stroke.UserID, stroke)
	return StrokesCollection.Insert(stroke)
}

func (p *Persistor) findNear(stroke *Stroke) ([]Stroke, error) {
	results := []Stroke{}
	query := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": stroke.Location,
				},
				"$maxDistance": maxDistance,
			},
		},
		"user_id": bson.M{
			"$ne": stroke.UserID,
		},
		"created_at": bson.M{
			"$gte": time.Now().Add(-1 * secondsRange * time.Second),
			"$lte": time.Now().Add(secondsRange * time.Second),
		},
	}
	err := StrokesCollection.Find(query).All(&results)
	utils.Log.Infof("Query executed by %s: %v", stroke.UserID, query)
	utils.Log.Infof("Actor %s found matches %d", stroke.UserID, len(results))

	return results, err
}
