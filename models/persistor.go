package models

import (
	"time"

	"github.com/mrkaspa/geoserver/utils"
	"gopkg.in/mgo.v2/bson"
)

type Persistance interface {
	PersistAndFind(StrokeNear) ([]Stroke, error)
	Persist(Stroke) error
	FindStrokes(string) ([]Stroke, error)
}

// Persistor takes care of persisting strokes and finding near ones
type persistor struct {
}

func NewPersistor() Persistance {
	return persistor{}
}

func (p persistor) PersistAndFind(sn StrokeNear) ([]Stroke, error) {
	utils.Log.Infof("Persistor executing PersistAndFind: %v", sn)
	if err := p.save(sn.Stroke); err != nil {
		return nil, err
	}

	nearUsers, err := p.findNear(sn)
	return nearUsers, err
}

func (p persistor) Persist(s Stroke) error {
	utils.Log.Infof("Persistor executing Persist: %s", s.UserID)
	return p.save(s)
}

func (p persistor) FindStrokes(username string) ([]Stroke, error) {
	strokes, err := p.history(username)
	return strokes, err
}

func (p *persistor) save(stroke Stroke) error {
	stroke.CreatedAt = time.Now()
	utils.Log.Infof("Persisting %s stroke: %v", stroke.UserID, stroke)
	return StrokesCollection.Insert(stroke)
}

func (p *persistor) history(username string) ([]Stroke, error) {
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

func (p *persistor) findNear(strokeNear StrokeNear) ([]Stroke, error) {
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
