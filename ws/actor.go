package ws

import (
	"time"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

const timeAlive = 3

type actorStatus int

const (
	_                 = iota
	alive actorStatus = 1
	dead
)

// Represents an user in the system that is doing a cheers
type actor struct {
	name             string
	status           actorStatus
	info             []byte
	timer            *time.Timer
	connections      []*connection
	matchedActors    map[*actor]bool
	sentActors       map[*actor]bool
	addConnection    chan *connection
	removeConnection chan *connection
	strokes          chan *models.Stroke
	nearUsers        chan []models.Stroke
	responses        chan []byte
	ping             chan *actor
	pong             chan *actor
	poisonPill       chan bool
}

// inits an actor
func newActor(name string) *actor {
	utils.Log.Infof("Creating actor: %s", name)
	return &actor{
		name:             name,
		status:           alive,
		connections:      []*connection{},
		addConnection:    make(chan *connection),
		removeConnection: make(chan *connection),
		matchedActors:    make(map[*actor]bool),
		sentActors:       make(map[*actor]bool),
		strokes:          make(chan *models.Stroke),
		responses:        make(chan []byte),
		nearUsers:        make(chan []models.Stroke, 256),
		ping:             make(chan *actor, 256),
		pong:             make(chan *actor, 256),
		poisonPill:       make(chan bool, 1),
	}
}

// run the actor that represents the user
func (a *actor) run() {
	utils.Log.Infof("Running actor: %s", a.name)
	for {
		select {
		case <-a.poisonPill:
			a.die()
			return
		case conn := <-a.addConnection:
			utils.Log.Infof("adding connection: %s to: %s - len before add: %d", conn.name, a.name, len(a.connections))
			a.connections = append(a.connections, conn)
		case conn := <-a.removeConnection:
			a.removeConnectionBy(conn)
		case strokeVar := <-a.strokes:
			a.persist(strokeVar)
		case users := <-a.nearUsers:
			for _, u := range users {
				searchActorVar := searchActor{
					name:     u.UserID,
					response: make(chan *actor),
				}
				SearcherVar.search <- &searchActorVar
				actorRef := <-searchActorVar.response
				actorRef.ping <- a
				a.matchedActors[actorRef] = false
			}
		case actorPing := <-a.ping:
			if _, ok := a.matchedActors[actorPing]; !ok {
				utils.Log.Infof("%s received a PING from %s", a.name, actorPing.name)
				a.matchedActors[actorPing] = false
				a.broadcast(actorPing)
				actorPing.pong <- a
			}
		case actorPong := <-a.pong:
			utils.Log.Infof("%s received a PONG from %s", a.name, actorPong.name)
			utils.Log.Infof("Dead of %s depends on %s", actorPong.name, a.name)
			a.matchedActors[actorPong] = true
			actorPong.timer.Stop()
			a.broadcast(actorPong)
		}
	}
}

// timeout to kill the actor started after the actor is registered
func (a *actor) startTimer() {
	a.timer = time.NewTimer(timeAlive * time.Second)
	<-a.timer.C
	a.poisonPill <- true
}

// sends the persist message
func (a *actor) persist(strokeVar *models.Stroke) {
	persistorVar := models.NewPersistor()
	persistorVar.UsersFound = a.nearUsers
	a.info = []byte(strokeVar.Info)
	persistorVar.PersistAndFind <- strokeVar
}

// sends data to all connectios
func (a *actor) broadcast(actorMatched *actor) {
	ok := a.sentActors[actorMatched]
	sent := false
	if !ok {
		sent = true
		a.sentActors[actorMatched] = true
		for _, conn := range a.connections {
			utils.Log.Infof("Sending to the connection %s the info %s of the actor %s", a.name, string(actorMatched.info), actorMatched.name)
			conn.send <- actorMatched.info
		}
	}
	utils.Log.Infof("Sent? %t, %s broadcasting %s, found: %t, on connections: %d", sent, a.name, actorMatched.name, ok, len(a.connections))
}

// removes a connection
func (a *actor) removeConnectionBy(conn *connection) {
	for i, c := range a.connections {
		if c == conn {
			a.connections = append(a.connections[:i], a.connections[i+1:]...)
			return
		}
	}
}

// finish an actor
func (a *actor) die() {
	// kills the referenced actors
	if a.status == alive {
		utils.Log.Infof("Actor dying: %s -- with %d connections", a.name, len(a.connections))

		a.status = dead

		//kills ponged actors
		for actorRef, ponged := range a.matchedActors {
			if ponged {
				actorRef.poisonPill <- true
			}
		}
		// closes all the actor connections
		for _, conn := range a.connections {
			utils.Log.Infof("Closing connection: %s", conn.name)
			conn.poisonPill <- true
			a.removeConnection <- conn
		}
		SearcherVar.unregister <- a.name
	}
}
