package handler

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("searcher", func() {

	actorName := "demo"

	BeforeEach(func() {
		SearcherVar.Clean <- true
	})

	It("should register the actor", func() {
		response := make(chan *actor)
		SearcherVar.register <- &registerActorWithResponse{
			name:     actorName,
			response: response,
		}
		actor := <-response
		Expect(actor.name).To(Equal(actorName))
	})

	Context("with one actor registered", func() {

		var createdActor *actor

		BeforeEach(func() {
			response := make(chan *actor)
			SearcherVar.register <- &registerActorWithResponse{
				name:     actorName,
				response: response,
			}
			createdActor = <-response
		})

		It("should find the actor", func() {
			response := make(chan *actor)
			SearcherVar.search <- &searchActorWithResponse{
				name:     actorName,
				response: response,
			}
			actorFound := <-response
			Expect(actorFound).To(Equal(createdActor))
		})

		It("should unregister the actor when this dies", func() {
			SearcherVar.unregister <- actorName
			response := make(chan *actor)
			SearcherVar.search <- &searchActorWithResponse{
				name:     actorName,
				response: response,
			}
			actorFound := <-response
			Expect(actorFound).To(BeNil())
		})

	})

})
