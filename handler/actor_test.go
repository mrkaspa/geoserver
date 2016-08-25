package handler

import (
	"time"

	"github.com/mrkaspa/geoserver/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("actor", func() {

	actorName := "demo"

	It("should die after when receives the poison pill seconds", func() {
		createdActor := newActor(actorName, models.NewMockPersistor)
		<-createdActor.poisonPill
		Expect(createdActor.status).To(Equal(dead))
	})

	It("should be unregistered after dying", func() {
		response := make(chan *actor)
		SearcherVar.register <- &registerActorWithResponse{
			name:     actorName,
			response: response,
		}
		createdActor := <-response
		createdActor.poisonPill <- true
		time.Sleep(1 * time.Second)

		SearcherVar.search <- &searchActorWithResponse{
			name:     actorName,
			response: response,
		}
		actorFound := <-response
		Expect(actorFound).To(BeNil())
	})

})
