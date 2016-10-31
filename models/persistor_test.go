package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/mrkaspa/geoserver/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	utils.LoadEnv("../.env_test")
	utils.InitLogger()
	InitDB()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestPersistor_Persist(t *testing.T) {
	cleanDB()
	persistor := NewPersistor()
	err := persistor.Persist(Stroke{
		UserID:   "a1",
		Info:     bson.M{},
		Location: []float64{-79.38066843, 43.65483486},
	})
	assert.Nil(t, err)
	strokes, err := persistor.FindStrokes("a1")
	assert.Nil(t, err)
	assert.NotEmpty(t, strokes)
}

func TestPersistor_PersistAndFind(t *testing.T) {
	cleanDB()
	a1 := Stroke{
		UserID:   "a1",
		Info:     bson.M{},
		Location: []float64{-79.38066843, 43.65483486},
	}

	a2 := Stroke{
		UserID:   "a2",
		Info:     bson.M{},
		Location: []float64{-79.38066843, 43.65483486},
	}

	a3 := Stroke{
		UserID:   "a2",
		Info:     bson.M{},
		Location: []float64{59.38066843, 43.65483486},
	}

	a4 := Stroke{
		UserID:   "a4",
		Info:     bson.M{"type": "demo"},
		Location: []float64{59.38066843, 43.65483486},
	}

	a5 := Stroke{
		UserID:   "a5",
		Info:     bson.M{},
		Location: []float64{59.38066843, 43.65483486},
	}

	params := map[string]interface{}{"type": "demo"}

	strokesTests := []struct {
		stroke          Stroke
		strokeNear      StrokeNear
		expectedResults int
		savedStrokes int
	}{
		{a1, StrokeNear{5, 10, a2, true, nil}, 1, 2},
		{a1, StrokeNear{5, 10, a1, true, nil}, 0, 2},
		{a1, StrokeNear{5, 10, a3, true, nil}, 0, 2},
		{a4, StrokeNear{5, 10, a5, true, params}, 1, 2},
		{a4, StrokeNear{5, 10, a5, false, params}, 1, 1},
	}

	for i, tt := range strokesTests {
		cleanDB()
		fmt.Printf("sample %d\n", i)
		persistor := NewPersistor()

		err := persistor.Persist(tt.stroke)
		assert.Nil(t, err)

		strokes, err := persistor.PersistAndFind(tt.strokeNear)
		assert.Nil(t, err)
		assert.Equal(t, len(strokes), tt.expectedResults)

		results := []Stroke{}
		err = StrokesCollection.Find(bson.M{}).All(&results)
		assert.Nil(t, err)
		assert.Equal(t, len(results), tt.savedStrokes)
	}
}

func cleanDB() {
	StrokesCollection.RemoveAll(bson.M{})
}
