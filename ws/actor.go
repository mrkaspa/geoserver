package ws

import (
	"time"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

const defaultLifeTime = 3

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
	lifeTime         time.Time
	connections      []*connection
	matchedActors    map[*actor]bool
	sentActors       map[*actor]bool
	addConnection    chan *connection
	removeConnection chan *connection
	strokesNear      chan *models.StrokeNear
	nearUsers        chan []models.Stroke
	responses        chan []byte
	ping             chan *actor
	pong             chan *actor
	poisonPill       chan bool
}

// inits an actor
func newActor(name string) *actor {
	utils.Log.Infof("Creating actor: %s", name)
	actor := &actor{
		name:             name,
		status:           alive,
		connections:      []*connection{},
		lifeTime:         time.Now(),
		addConnection:    make(chan *connection),
		removeConnection: make(chan *connection),
		matchedActors:    make(map[*actor]bool),
		sentActors:       make(map[*actor]bool),
		strokesNear:      make(chan *models.StrokeNear),
		responses:        make(chan []byte),
		nearUsers:        make(chan []models.Stroke, 256),
		ping:             make(chan *actor, 256),
		pong:             make(chan *actor, 256),
		poisonPill:       make(chan bool, 1),
	}
	go actor.startTimer(defaultLifeTime)
	go actor.run()
	return actor
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
		case strokeVar := <-a.strokesNear:
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
				diffTime := actorPing.lifeTime.Sub(a.lifeTime)
				a.dieLater(int(diffTime.Seconds()))
				actorPing.pong <- a
			}
		case actorPong := <-a.pong:
			utils.Log.Infof("%s received a PONG from %s", a.name, actorPong.name)
			a.matchedActors[actorPong] = true
			a.broadcast(actorPong)
		}
	}
}

// timeout to kill the actor started after the actor is registered
func (a *actor) startTimer(seconds int) {
	a.lifeTime = a.lifeTime.Add(time.Duration(seconds) * time.Second)
	a.timer = time.NewTimer(time.Duration(seconds) * time.Second)
	<-a.timer.C
	a.poisonPill <- true
}

// postpones death
func (a *actor) dieLater(seconds int) {
	if seconds > 0 {
		a.timer.Stop()
		utils.Log.Infof("Actor %s die later %d", a.name, seconds)
		go a.startTimer(seconds)
	}
}

// sends the persist message
func (a *actor) persist(strokeNear *models.StrokeNear) {
	persistor := models.NewPersistor()
	persistor.UsersFound = a.nearUsers
	a.info = []byte(strokeNear.Stroke.Info)
	persistor.PersistAndFind <- strokeNear
	if <-persistor.Saved {
		go a.dieLater(strokeNear.TimeRange)
	}
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
		utils.Log.Infof("Actor dying: %s -- with %d connections -- should die %v", a.name, len(a.connections), a.lifeTime)

		a.status = dead

		// closes all the actor connections
		for _, conn := range a.connections {
			utils.Log.Infof("Closing connection: %s", conn.name)
			conn.poisonPill <- true
			a.removeConnection <- conn
		}
		SearcherVar.unregister <- a.name
	}
}
