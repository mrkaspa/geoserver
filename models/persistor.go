package models

import (
	"time"

	"github.com/mrkaspa/geoserver/utils"
	"gopkg.in/mgo.v2/bson"
)

type Persistance interface {
	PersistAndFind() chan PersistAndFindWithResponse
	Persist() chan PersistWithResponse
	FindStrokes() chan FindStrokesWithResponse
}

// Persistor takes care of persisting strokes and finding near ones
type Persistor struct {
	persistAndFind chan PersistAndFindWithResponse
	persist        chan PersistWithResponse
	findStrokes    chan FindStrokesWithResponse
}

func NewPersistor() Persistance {
	persistor := Persistor{
		persistAndFind: make(chan PersistAndFindWithResponse, 256),
		persist:        make(chan PersistWithResponse, 256),
		findStrokes:    make(chan FindStrokesWithResponse, 256),
	}
	go persistor.run()
	return persistor
}

func (p Persistor) PersistAndFind() chan PersistAndFindWithResponse {
	return p.persistAndFind
}

func (p Persistor) Persist() chan PersistWithResponse {
	return p.persist
}

func (p Persistor) FindStrokes() chan FindStrokesWithResponse {
	return p.findStrokes
}

func (p Persistor) run() {
	defer func() {
		if r := recover(); r != nil {
			utils.Log.Infof("Recovered in persistor.Run()")
		}
	}()
	for {
		select {
		case pwr := <-p.persist:
			stroke := pwr.Stroke
			utils.Log.Infof("Persistor executing Persist: %s", stroke.UserID)
			if err := p.save(stroke); err != nil {
				pwr.Response <- false
				continue
			}
			pwr.Response <- true
		case pfwr := <-p.persistAndFind:
			strokeNear := pfwr.StrokeNear
			utils.Log.Infof("Persistor executing PersistAndFind: %v", strokeNear)
			if err := p.save(strokeNear.Stroke); err != nil {
				pfwr.SavedResponse <- false
				continue
			}

			pfwr.SavedResponse <- true
			nearUsers, err := p.findNear(strokeNear)
			if err != nil {
				continue
			}
			pfwr.UsersResponse <- nearUsers
		case fs := <-p.findStrokes:
			strokes, err := p.history(fs.Username)
			if err != nil {
				continue
			}
			fs.Response <- strokes
		}
	}
}

func (p *Persistor) save(stroke Stroke) error {
	stroke.CreatedAt = time.Now()
	utils.Log.Infof("Persisting %s stroke: %v", stroke.UserID, stroke)
	return StrokesCollection.Insert(stroke)
}

func (p *Persistor) history(username string) ([]Stroke, error) {
	results := []Stroke{}
	query := buildHistoryQuery(username)
	err := StrokesCollection.Find(query).Sort("-created_at").All(&results)
	utils.Log.Infof("Query executed by %s", username)
	utils.Log.Infof("Actor %s found matches %d", username, len(results))
	return results, err
}

func buildHistoryQuery(username string) bson.M {
	return bson.M{"user_id": username}
}

func (p *Persistor) findNear(strokeNear StrokeNear) ([]Stroke, error) {
	results := []Stroke{}
	query := buildNearQuery(strokeNear)
	err := StrokesCollection.Find(query).Sort("-created_at").All(&results)
	utils.Log.Infof("Query executed by %s: %v", strokeNear.Stroke.UserID, query)
	utils.Log.Infof("Actor %s found matches %d", strokeNear.Stroke.UserID, len(results))
	return results, err
}

func buildNearQuery(strokeNear StrokeNear) bson.M {
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
