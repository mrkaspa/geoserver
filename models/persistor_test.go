package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("persistor", func() {

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

	clearDB := func() {
		StrokesCollection.RemoveAll(bson.M{})
	}

	AfterEach(func() {
		clearDB()
	})

	It("should save the stroke", func() {
		persistor := NewPersistor()
		persistor.Persist() <- &Stroke{
			UserID:   "a1",
			Info:     "a1",
			Location: []float64{-79.38066843, 43.65483486},
		}
		saved := <-persistor.Saved()
		Expect(saved).To(BeTrue())
	})

	It("should get the expected matches", func() {
		persistor := NewPersistor()
		for _, tt := range strokesTests {
			persistor.Persist() <- &tt.stroke
			Expect(<-persistor.Saved()).To(BeTrue())
			persistor.PersistAndFind() <- &tt.strokeNear
			Expect(<-persistor.Saved()).To(BeTrue())
			matches := <-persistor.UsersFound()
			Expect(len(matches)).To(Equal(tt.expectedResults))
			clearDB()
		}
	})

})
