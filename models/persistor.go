package models

import (
	"time"

	"github.com/mrkaspa/geoserver/utils"
	"gopkg.in/mgo.v2/bson"
)

// Persistor takes care of persisting strokes and finding near ones
type Persistor struct {
	PersistAndFind chan *StrokeNear
	Persist        chan *Stroke
	UsersFound     chan []Stroke
	Saved          chan bool
}

func NewPersistor() *Persistor {
	persistor := &Persistor{
		PersistAndFind: make(chan *StrokeNear, 1),
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
		p.save(stroke)
	case strokeNear := <-p.PersistAndFind:
		utils.Log.Infof("Persistor executing PersistAndFind: %s", strokeNear)
		if p.save(&strokeNear.Stroke) {
			nearUsers, err := p.findNear(strokeNear)
			if err != nil {
				panic(err)
			}
			p.UsersFound <- nearUsers
		}
	}
}

func (p *Persistor) save(stroke *Stroke) bool {
	stroke.CreatedAt = time.Now()
	utils.Log.Infof("Persisting %s stroke: %v", stroke.UserID, stroke)
	err := StrokesCollection.Insert(stroke)
	if err != nil {
		p.Saved <- false
		return false
	}
	p.Saved <- true
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
