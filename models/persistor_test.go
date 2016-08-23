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
		response := make(chan bool)
		persistor := NewPersistor()
		persistor.Persist() <- PersistWithResponse{
			Stroke: Stroke{
				UserID:   "a1",
				Info:     "a1",
				Location: []float64{-79.38066843, 43.65483486},
			},
			Response: response,
		}
		Expect(<-response).To(BeTrue())

		strokesResponse := make(chan []Stroke)
		persistor.FindStrokes() <- FindStrokesWithResponse{
			Username: "a1",
			Response: strokesResponse,
		}

		Expect(<-strokesResponse).ToNot(BeEmpty())
	})

	for i, tt := range strokesTests {
		It("should get the expected matches", func() {
			persistor := NewPersistor()
			fmt.Printf("running %d", i)
			savedResponse := make(chan bool)
			usersResponse := make(chan []Stroke)

			persistor.Persist() <- PersistWithResponse{
				Stroke:   tt.stroke,
				Response: savedResponse,
			}
			Expect(<-savedResponse).To(BeTrue())

			persistor.PersistAndFind() <- PersistAndFindWithResponse{
				StrokeNear:    tt.strokeNear,
				SavedResponse: savedResponse,
				UsersResponse: usersResponse,
			}
			Expect(<-savedResponse).To(BeTrue())
			matches := <-usersResponse
			Expect(len(matches)).To(Equal(tt.expectedResults))
		})
	}

})
