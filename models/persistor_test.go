package models

import (
	"fmt"

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

	cleanDB := func() {
		StrokesCollection.RemoveAll(bson.M{})
	}

	BeforeEach(func() {
		cleanDB()
	})

	It("should save the stroke", func() {
		persistor := NewPersistor()
		err := persistor.Persist(Stroke{
			UserID:   "a1",
			Info:     "a1",
			Location: []float64{-79.38066843, 43.65483486},
		})
		Expect(err).To(BeNil())

		strokes, err := persistor.FindStrokes("a1")
		Expect(err).To(BeNil())
		Expect(strokes).ToNot(BeEmpty())
	})

	for i, tt := range strokesTests {
		It("should get the expected matches", func() {
			persistor := NewPersistor()
			fmt.Printf("running %d", i)

			err := persistor.Persist(tt.stroke)
			Expect(err).To(BeNil())

			strokes, err := persistor.PersistAndFind(tt.strokeNear)
			Expect(err).To(BeNil())
			Expect(len(strokes)).To(Equal(tt.expectedResults))
		})
	}

})
