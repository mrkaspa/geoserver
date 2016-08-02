package handler

import (
	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

type searcher struct {
	persistorCreator func() models.Persistance
	directory        map[string]*actor
	search           chan *searchActorWithResponse
	register         chan *registerActorWithResponse
	unregister       chan string
	Clean            chan bool
}

// Searcher for the actors on the system
var SearcherVar *searcher

func InitSearcher(persistorCreator func() models.Persistance) {
	// InitSearcher run it
	SearcherVar = &searcher{
		persistorCreator: persistorCreator,
		directory:        make(map[string]*actor),
		search:           make(chan *searchActorWithResponse, 256),
		register:         make(chan *registerActorWithResponse, 256),
		unregister:       make(chan string, 256),
		Clean:            make(chan bool),
	}
	go SearcherVar.Run()
}

// Run the searcher
func (s *searcher) Run() {
	defer func() {
		if r := recover(); r != nil {
			utils.Log.Infof("Recovered in SearcherVar.Run()")
		}
	}()

	for {
		select {
		case search := <-s.search:
			actorRef := s.directory[search.name]
			search.response <- actorRef
		case register := <-s.register:
			// creates or find an actor
			actorRef, ok := s.directory[register.name]
			utils.Log.Infof("Looking for actor: %s --- %v", register.name, actorRef)
			if !ok {
				actorRef = newActor(register.name, s.persistorCreator)
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
