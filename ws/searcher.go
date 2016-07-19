package ws

import "github.com/mrkaspa/geoserver/utils"

type searcher struct {
	directory  map[string]*actor
	search     chan *searchActor
	register   chan *registerActor
	unregister chan string
	Clean      chan bool
}

// Searcher for the actors on the system
var SearcherVar *searcher

// Run the searcher
func (s *searcher) Run() {
	for {
		select {
		case search := <-s.search:
			actorRef, _ := s.directory[search.name]
			search.response <- actorRef
		case register := <-s.register:
			// creates or find an actor
			actorRef, ok := s.directory[register.name]
			utils.Log.Infof("Looking for actor: %s --- %v", register.name, actorRef)
			if !ok {
				actorRef = newActor(register.name)
				s.directory[register.name] = actorRef
			} else {
				utils.Log.Infof("Actor found: %s", register.name)
			}
			register.response <- actorRef
		case username := <-s.unregister:
			if _, ok := s.directory[username]; ok {
				delete(s.directory, username)
			}
		case <-s.Clean:
			s.directory = make(map[string]*actor)
		}
	}
}
