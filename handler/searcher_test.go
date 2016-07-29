package handler

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"
)

var _ = Describe("searcher", func() {

	actorName := "demo"

	BeforeEach(func() {
		SearcherVar.Clean <- true
	})

	It("should register the actor", func() {
		Expect(SearcherVar.directory).To(BeEmpty())
		response := make(chan *actor)
		SearcherVar.register <- &registerActor{
			name:     actorName,
			response: response,
		}
		actor := <-response
		Expect(actor.name).To(Equal(actorName))
		Expect(SearcherVar.directory).ToNot(BeEmpty())
	})

	Context("with one actor registered", func() {

		var createdActor *actor

		BeforeEach(func() {
			response := make(chan *actor)
			SearcherVar.register <- &registerActor{
				name:     actorName,
				response: response,
			}
			createdActor = <-response
		})

		It("should find the actor", func() {
			response := make(chan *actor)
			SearcherVar.search <- &searchActor{
				name:     actorName,
				response: response,
			}
			actorFound := <-response
			Expect(actorFound).To(Equal(createdActor))
		})

		It("should unregister the actor when this dies", func() {
			SearcherVar.unregister <- actorName
			time.Sleep(1 * time.Second)
			Expect(SearcherVar.directory).To(BeEmpty())
		})

	})

})
