package models

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"os"
	"github.com/mrkaspa/geoserver/utils"
	"github.com/stretchr/testify/assert"
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
		Info:     "a1",
		Location: []float64{-79.38066843, 43.65483486},
	})
	assert.Nil(t, err)
	strokes, err := persistor.FindStrokes("a1")
	assert.Nil(t, err)
	assert.NotEmpty(t, strokes)
}

func TestPersistor_FindStrokes(t *testing.T) {
	cleanDB()
	a1 := Stroke{
		UserID:   "a1",
		Info:     "a1",
		Location: []float64{-79.38066843, 43.65483486},
	}

	a2 := Stroke{
		UserID:   "a2",
		Info:     "a2",
		Location: []float64{-79.38066843, 43.65483486},
	}

	a3 := Stroke{
		UserID:   "a2",
		Info:     "a2",
		Location: []float64{59.38066843, 43.65483486},
	}

	strokesTests := []struct {
		stroke          Stroke
		strokeNear      StrokeNear
		expectedResults int
	}{
		{a1, StrokeNear{5, 10, a2}, 1},
		{a1, StrokeNear{5, 10, a1}, 0},
		{a1, StrokeNear{5, 10, a3}, 0},
	}

	for _, tt := range strokesTests {
		cleanDB()
		persistor := NewPersistor()

		err := persistor.Persist(tt.stroke)
		assert.Nil(t, err)

		strokes, err := persistor.PersistAndFind(tt.strokeNear)
		assert.Nil(t, err)
		assert.Equal(t, len(strokes), tt.expectedResults)
	}
}

func cleanDB() {
	StrokesCollection.RemoveAll(bson.M{})
}